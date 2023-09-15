package v1

import (
	"context"
	"fio_service/internal/models"
	"fio_service/pkg/api/person"
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/proto"
	"io"
	"net/http"
	"strconv"
)

func (h *Handler) initPersonRoutes(api *gin.RouterGroup) {
	g := api.Group("/person")
	{
		g.POST("/", h.create)
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
// @Param			struct	body		models.Person	true	"Person"
// @Success		201		{object}	Resposne
// @Failure		400		{object}	Resposne
// @Failure		500		{object}	Resposne
// @Router			/person/create [put]
func (h *Handler) create(ctx *gin.Context) {
	var inp person.Person

	data, _ := io.ReadAll(ctx.Request.Body)

	err := proto.Unmarshal(data, &inp)
	if err != nil {
		newResponse(ctx, http.StatusBadRequest, "Incorrect input data format")
		return
	}

	if err := h.service.Person.Create(context.Background(), &models.Person{
		Name:        inp.Name,
		Surname:     inp.Surname,
		Patronymic:  inp.Patronymic,
		Age:         inp.Age,
		Gender:      models.PersonGender(inp.Gender),
		Nationality: inp.Nationality,
	}); err != nil {
		newResponse(ctx, http.StatusInternalServerError, "Can't create a person")
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
		newResponse(ctx, http.StatusBadRequest, "Incorrect person ID")
		return
	}
	p, err := h.service.Person.Get(context.Background(), uint64(id))
	if err != nil {
		newResponse(ctx, http.StatusInternalServerError, "Can't get a person")
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
		newResponse(ctx, http.StatusBadRequest, "Incorrect person ID")
		return
	}

	if err := h.service.Person.Delete(context.Background(), uint64(id)); err != nil {
		newResponse(ctx, http.StatusInternalServerError, "Can't delete a person")
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
	var inp person.Person
	param := ctx.Param("id")

	id, err := strconv.Atoi(param)
	if err != nil {
		newResponse(ctx, http.StatusBadRequest, "Incorrect person ID")
		return
	}

	data, _ := io.ReadAll(ctx.Request.Body)

	if err := proto.Unmarshal(data, &inp); err != nil {
		newResponse(ctx, http.StatusBadRequest, "Incorrect input data format")
		return
	}

	fields := make(models.PersonFieldsToUpdate)

	if inp.Name != "" {
		fields[models.PersonFieldName] = inp.Name
	}
	if inp.Surname != "" {
		fields[models.PersonFieldSurname] = inp.Surname
	}
	if inp.Patronymic != "" {
		fields[models.PersonFieldPatronymic] = inp.Patronymic
	}
	if inp.Age != 0 {
		fields[models.PersonFieldAge] = inp.Age
	}
	if inp.Gender != "" {
		fields[models.PersonFieldGender] = inp.Gender
	}
	if inp.Nationality != "" {
		fields[models.PersonFieldNationality] = inp.Nationality
	}

	if err := h.service.Person.Update(context.Background(), uint64(id), fields); err != nil {
		newResponse(ctx, http.StatusInternalServerError, "Can't update a person")
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

func (h *Handler) consumeResponseMessages() {
	err := h.service.Kafka.ConsumeMessages("response", h.handleResponseMessage)
	if err != nil {
		h.logger.Println(err)
	}
}

func (h *Handler) handleResponseMessage(message string) {
	h.responseCh <- []byte(message)
}
