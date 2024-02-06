package routes

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/sixfwa/fiber-api/database"
	"github.com/sixfwa/fiber-api/models"
	"gorm.io/datatypes"
)

type QuestionSerializer struct {
	ID            uint           `json:"id" gorm:"primaryKey"`
	Test          TestSerializer `gorm:"foreignKey:TestRefer"`
	CorrectAnswer string         `json:"correct_answer"`
	Answers       datatypes.JSON `json:"answers"`
	QuestionTitle string         `json:"title_question"`
	Item          ItemSerializer `gorm:"foreignKey:ItemRefer"`
	QuestionInfo  string         `json:"question_info"`
}

func CreateResponseQuestion(question models.Question, test TestSerializer, item ItemSerializer) QuestionSerializer {
	return QuestionSerializer{
		ID:            question.ID,
		Test:          test,
		CorrectAnswer: question.CorrectAnswer,
		Answers:       question.Answers,
		QuestionTitle: question.QuestionTitle,
		Item:          item,
		QuestionInfo:  question.QuestionInfo,
	}
}

func CreateQuestion(c *fiber.Ctx) error {
	var question models.Question

	if err := c.BodyParser(&question); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	var test models.Test
	if err := findTest(question.TestRefer, &test); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	var item models.Item
	if err := findItem(question.ItemRefer, &item); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	var year models.Year
	if err := findYear(item.YaerRefer, &year); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	// Обновить связанные объекты test и item
	database.Database.Db.Model(&test).Updates(&test)
	database.Database.Db.Model(&item).Updates(&item)

	database.Database.Db.Create(&question)

	responseTest := CreateResponseTest(test)
	responseYear := CreateResponseYear(year)
	responseItem := CreateResponseItem(item, responseYear)

	responseQuestion := CreateResponseQuestion(question, responseTest, responseItem)
	return c.Status(200).JSON(responseQuestion)
}

func GetQuestions(c *fiber.Ctx) error {
	questions := []models.Question{}

	database.Database.Db.Preload("Test").Preload("Item.Year").Find(&questions)

	responseQuestions := []QuestionSerializer{}
	for _, question := range questions {
		responseItem := CreateResponseItem(question.Item, CreateResponseYear(question.Item.Year))
		responseTest := CreateResponseTest(question.Test)

		responseQuestion := CreateResponseQuestion(question, responseTest, responseItem)

		responseQuestions = append(responseQuestions, responseQuestion)
	}
	return c.Status(200).JSON(responseQuestions)
}

func FindQuestion(id int, question *models.Question) error {
	database.Database.Db.Find(&question, "id = ?", id)
	if question.ID == 0 {
		return errors.New("order does not exist")
	}

	return nil
}
func GetQuestion(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON("Please ensure that :id is an integer")
	}

	var question models.Question
	if err := FindQuestion(id, &question); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	// Используйте Preload для предварительной загрузки связанных объектов
	database.Database.Db.Preload("Test").Preload("Item.Year").Find(&question)

	responseItem := CreateResponseItem(question.Item, CreateResponseYear(question.Item.Year))
	responseTest := CreateResponseTest(question.Test)

	responseQuestion := CreateResponseQuestion(question, responseTest, responseItem)
	return c.Status(200).JSON(responseQuestion)
}

func UpdateQuestion(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON("Please ensure that :id is an integer")
	}

	var question models.Question
	if err := FindQuestion(id, &question); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	type UpdateQuestionData struct {
		Test          models.Test    `json:"test"`
		CorrectAnswer string         `json:"correct_answer"`
		Answers       datatypes.JSON `json:"answers"`
		QuestionTitle string         `json:"title_question"`
		Item          models.Item    `json:"item"`
		QuestionInfo  string         `json:"question_info"`
	}

	var updateData UpdateQuestionData
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(500).JSON(err.Error())
	}

	// Используйте Preload для предварительной загрузки связанных объектов
	database.Database.Db.Preload("Item.Year").Find(&updateData.Item, "id = ?", updateData.Item.ID)

	question.Test = updateData.Test
	question.CorrectAnswer = updateData.CorrectAnswer
	question.Answers = updateData.Answers
	question.QuestionTitle = updateData.QuestionTitle
	question.Item = updateData.Item
	question.QuestionInfo = updateData.QuestionInfo

	database.Database.Db.Save(&question)
	database.Database.Db.Preload("Test").Preload("Item.Year").Find(&question)

	responseItem := CreateResponseItem(question.Item, CreateResponseYear(question.Item.Year))
	responseTest := CreateResponseTest(question.Test)

	responseQuestion := CreateResponseQuestion(question, responseTest, responseItem)
	return c.Status(200).JSON(responseQuestion)
}

func DeleteQuestion(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON("Please ensure that :id is an integer")
	}

	var question models.Question
	if err := FindQuestion(id, &question); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	// Удаление связанных объектов Item и Test
	database.Database.Db.Delete(&question.Item)
	database.Database.Db.Delete(&question.Test)

	if err := database.Database.Db.Delete(&question).Error; err != nil {
		return c.Status(404).JSON(err.Error())
	}

	return c.Status(200).SendString("Successfully Deleted Question")
}
