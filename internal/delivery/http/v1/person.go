package v1

import (
	"context"
	"encoding/json"
	"fio_service/internal/models"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"strconv"
)

// @host		localhost:8000
// @BasePath	/api/v1
func (h *Handler) initPersonRoutes(api *gin.RouterGroup) {
	g := api.Group("/person")
	{
		g.POST("/create", h.create)
		g.GET("/:id", h.get)
		g.DELETE("/:id", h.delete)
		g.PUT("/:id", h.update)
		g.GET("/list", h.getList)
	}
}

// @Summary		Create new Person
// @Tags			Person
// @Description	Create new Person
// @ModuleID		create
// @Accept			json
// @Produce		json
// @Param			struct	body		person.Person	true	"Person"
// @Success		201		{object}	Resposne
// @Failure		400		{object}	Resposne
// @Failure		500		{object}	Resposne
// @Router			/person/create [post]
func (h *Handler) create(ctx *gin.Context) {
	var p models.Person

	data, _ := io.ReadAll(ctx.Request.Body)

	err := json.Unmarshal(data, &p)
	if err != nil {
		newResponse(ctx, http.StatusBadRequest, "Incorrect input data format: "+err.Error())
		return
	}

	if err := h.service.Person.Create(context.Background(), &models.Person{
		Name:        p.Name,
		Surname:     p.Surname,
		Patronymic:  p.Patronymic,
		Age:         p.Age,
		Gender:      p.Gender,
		Nationality: p.Nationality,
	}); err != nil {
		newResponse(ctx, http.StatusInternalServerError, "Can't create a person: "+err.Error())
		return
	}

	ctx.JSON(http.StatusCreated, Resposne{"The person was successfully created"})
}

// @Summary		Get Person by ID
// @Tags			Person
// @Description	Get Person by ID
// @ModuleID		get
// @Accept			json
// @Produce		json
// @Param			id	path		integer	true	"person id"
// @Success		200	{object}	models.Person
// @Failure		400	{object}	Resposne
// @Failure		500	{object}	Resposne
// @Router			/person/{id} [get]
func (h *Handler) get(ctx *gin.Context) {
	param := ctx.Param("id")
	id, err := strconv.Atoi(param)
	if err != nil {
		newResponse(ctx, http.StatusBadRequest, "Incorrect person ID: "+err.Error())
		return
	}
	p, err := h.service.Person.Get(context.Background(), uint64(id))
	if err != nil {
		newResponse(ctx, http.StatusInternalServerError, "Can't get a person: "+err.Error())
		return
	}

	ctx.JSON(http.StatusOK, p)
}

// @Summary		Delete Person
// @Tags			Person
// @Description	Delete Person
// @ModuleID		delete
// @Accept			json
// @Produce		json
// @Param			id	path		integer	true	"person id"
// @Success		200	{object}	Resposne
// @Failure		400	{object}	Resposne
// @Failure		500	{object}	Resposne
// @Router			/person/{id} [delete]
func (h *Handler) delete(ctx *gin.Context) {
	param := ctx.Param("id")

	id, err := strconv.Atoi(param)
	if err != nil {
		newResponse(ctx, http.StatusBadRequest, "Incorrect person ID: "+err.Error())
		return
	}

	if err := h.service.Person.Delete(context.Background(), uint64(id)); err != nil {
		newResponse(ctx, http.StatusInternalServerError, "Can't delete a person:"+err.Error())
		return
	}

	ctx.JSON(http.StatusOK, Resposne{"Person was successfully deleted"})
}

// @Summary		Update Person
// @Tags			Person
// @Description	Update Person
// @Accept			json
// @Produce		json
// @ModuleID		update
// @Param			person	body		models.Person	true	"person update fields"
// @Param			id		path		integer			true	"person id"
// @Success		200		{object}	Resposne
// @Failure		400		{object}	Resposne
// @Failure		500		{object}	Resposne
// @Router			/person/{id} [put]
func (h *Handler) update(ctx *gin.Context) {
	var p models.Person
	param := ctx.Param("id")

	id, err := strconv.Atoi(param)
	if err != nil {
		newResponse(ctx, http.StatusBadRequest, "Incorrect person ID: "+err.Error())
		return
	}

	data, _ := io.ReadAll(ctx.Request.Body)

	if err := json.Unmarshal(data, &p); err != nil {
		newResponse(ctx, http.StatusBadRequest, "Incorrect input data format: "+err.Error())
		return
	}

	fields := make(models.PersonFieldsToUpdate)

	if p.Name != "" {
		fields[models.PersonFieldName] = p.Name
	}
	if p.Surname != "" {
		fields[models.PersonFieldSurname] = p.Surname
	}
	if p.Patronymic != "" {
		fields[models.PersonFieldPatronymic] = p.Patronymic
	}
	if p.Age != 0 {
		fields[models.PersonFieldAge] = p.Age
	}
	if p.Gender != "" {
		fields[models.PersonFieldGender] = p.Gender
	}
	if p.Nationality != "" {
		fields[models.PersonFieldNationality] = p.Nationality
	}

	if err := h.service.Person.Update(context.Background(), uint64(id), fields); err != nil {
		newResponse(ctx, http.StatusInternalServerError, "Can't update a person: "+err.Error())
		return
	}

	ctx.JSON(http.StatusOK, Resposne{"Person was successfully updated"})
}

// @Summary		Get Person List
// @Tags			Person
// @Description	Get Person List
// @ModuleID		get
// @Accept			json
// @Produce		json
// @Success		200	{object}	[]models.Person
// @Failure		400	{object}	Resposne
// @Failure		500	{object}	Resposne
// @Router			/person/list [get]
func (h *Handler) getList(ctx *gin.Context) {
	p, err := h.service.Person.GetList(context.Background())
	if err != nil {
		newResponse(ctx, http.StatusInternalServerError, "Can't get a person list")
		return
	}

	ctx.JSON(http.StatusOK, p)
}

func (h *Handler) consumeMessages() {
	err := h.service.Kafka.ConsumeMessages("FIO", h.handleMessage)
	if err != nil {
		h.logger.Println(err)
	}
}

func (h *Handler) handleMessage(message string) {
	h.logger.Info("message received: \n" + message)
	var p models.Person

	err := json.Unmarshal([]byte(message), &p)
	if err != nil {
		if errIn := h.service.Kafka.SendMessages("FIO_FAILED", "invalid format: "+err.Error()); errIn != nil {
			h.logger.Error("can't send error message to the topic")
		}
		return
	}

	if err := h.service.Person.CreateWithEnrichment(context.Background(), &models.Person{
		Name:       p.Name,
		Surname:    p.Surname,
		Patronymic: p.Patronymic,
	}); err != nil {
		if errIn := h.service.Kafka.SendMessages("FIO_FAILED", "can't create a person: "+err.Error()); errIn != nil {
			h.logger.Error("can't send error message to the topic")
		}
		return
	}
}
