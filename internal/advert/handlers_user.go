package advert

import (
    "fmt"
    "time"
    "errors"
    "github.com/gofiber/fiber/v2"
    "github.com/tvitcom/local-adverts/internal/config"
    "github.com/tvitcom/local-adverts/pkg/util"  
    
)

func (res resource) pageUserProfile(c *fiber.Ctx) error {
    uid := util.Pkeyer(c.Locals("iam"))
    user, err := res.agregator.GetUserById(c.UserContext(), uid)
    if err != nil {
        return err
    }
    return c.Render("userprofile", fiber.Map{
        "msg": "редактирование данных",
        "user": user,
    })
}

func (res resource) handlerUserProfile(c *fiber.Ctx) error {
    uid := util.Pkeyer(c.Locals("iam"))
    curruser, err := res.agregator.GetUserById(c.UserContext(), uid)
    if err != nil {
        return c.Status(412).Render("error", fiber.Map{"msg": err})
    }

    form := new(ProfileForm)
    if err := c.BodyParser(form); err != nil {
        res.logger.With(c.UserContext()).Error(err.Error())
        return c.Status(412).Render("error", fiber.Map{
            "msg": err,
            "user": curruser,
        })
    }

    form.Tel = util.PhoneNormalisation(form.Tel)   
    // Form validation    
    if err := form.Validate(); err != nil {
        return c.Status(412).Render("error", fiber.Map{"msg": err})
    }

    // Parse the multipart form:
    imagefname := ""
    if ff, err := c.MultipartForm(); err == nil {
        files := ff.File["picture"]// => []*multipart.FileHeader

        // Loop through files:
        for i, ff := range files {
            if i > 1 {
                return errors.New("Файлов слишком много для сайта")
            }
            // Start Image convey:      
            userpicfpath, err := util.MakeUploadDirByUserId(config.PictureUserPath, util.Stringer(uid))
            fmt.Println("TO pictureMultypartFile:", err)
            if err != nil {
                return c.Status(501).Render("error", fiber.Map{
                    "msg": err,
                })
            }
            imagefname = util.GetMD5Hash(fmt.Sprintf("%s", time.Now().UnixNano()))+".jpg"
            err = c.SaveFile(ff, config.UploadedPath + imagefname)
            if err != nil {
                return c.Status(500).Render("error", fiber.Map{
                    "msg": err,
                })
            }
            err = util.ImagefileValidations(config.UploadedPath + imagefname)
            if err != nil {
                
                return c.Status(501).Render("error", fiber.Map{
                    "msg": "Загруженная картинка не подходит для сайта",
                })
            }
            err = util.ImagefileResizing(config.UploadedPath + imagefname, userpicfpath + imagefname, 75)
            if err != nil {
                                return c.Status(500).Render("error", fiber.Map{
                    "msg": "Загруженная картинка не уменьшена",
                })
            }
            // Remove current users picture
            err = util.ImagefileRemove(userpicfpath + curruser.Picture)
            if err != nil {
               fmt.Println("Нельзя удалить прошлую картинку пользователя:", userpicfpath + curruser.Picture)
            }
            // Remove uploads temporary picture
            err = util.ImagefileRemove(config.UploadedPath + imagefname)
            if err != nil {
                return c.Status(500).Render("error", fiber.Map{
                    "msg": "Нельзя удалить временную картинку",
                })
            }
            err = util.ImagefileProgressiveOptimisation(c.UserContext(), userpicfpath + imagefname, "", false)
            if err != nil {
                c.Status(500)
                return c.Render("error", fiber.Map{
                    "msg": "Загруженная картинка не уменьшена",
                })
            }
        }
    }

    err = res.agregator.UpdateUser(c.UserContext(), form, uid, imagefname)
    if err != nil {
        return c.Status(412).Render("error", fiber.Map{"msg": err})
    }
    return c.Redirect("/my/userprofile.html")
}

func (res resource) pageUserAdverts(c *fiber.Ctx) error {
    uid := util.Pkeyer(c.Locals("iam"))
    curruser, err := res.agregator.GetUserById(c.UserContext(), uid)
    if err != nil {
        return c.Status(412).Render("error", fiber.Map{"msg": err})
    }
    advertsdisplay, err := res.agregator.GetAdvertsDisplayByUserId(c.UserContext(), uid)
    if err != nil {
        return c.Status(500).Redirect("/error.html?msg=Ошибка работы сайта")
    }

    return c.Render("useradverts", fiber.Map{
        "msg": "useradverts page: page: Coming soon!",
        "advertsdisplay": advertsdisplay,
        "user": curruser,
    })
}

func (res resource) handlerDeleteAdvert(c *fiber.Ctx) error {
    uid := util.Pkeyer(c.Locals("iam"))
    form := new(DeleteAdvertForm)
    if err := c.BodyParser(form); err != nil && uid != 0 {
        return c.Status(500).Redirect("/error.html?msg=Ошибка обработки формы для удаления объявления")
    }
    if err := res.agregator.DeleteAdvertsData(c.UserContext(), form.AdvertId); err != nil {
        res.logger.With(c.UserContext()).Error(err.Error())
        return c.Status(500).Redirect("/error.html?msg=Ошибка удаления объявления")
    }
    return c.Status(201).Render("thanks", fiber.Map{
        "msg": "Объявление и его фотографии успешно удалены",
    })
}

func (res resource) pageUserMessages(c *fiber.Ctx) error {
    uid := util.Pkeyer(c.Locals("iam"))
    curruser, err := res.agregator.GetUserById(c.UserContext(), uid)
    if err != nil {
        return c.Status(412).Render("error", fiber.Map{"msg": err})
    }
    senders, err := res.agregator.GetMessagesSendersByUserId(c.UserContext(), uid)
    if err != nil {
        return c.Status(412).Render("error", fiber.Map{"msg": err})
    }
    return c.Render("usermessages", fiber.Map{
        "msg": "usermessages page: page: Coming soon!",
        "senders": senders,
        "user": curruser,
    })
}

func (res resource) pageUserList(c *fiber.Ctx) error {
    uid := util.Pkeyer(c.Locals("iam"))
    if uid > 1 { // Non admin with id=1 go ahead
        return c.Status(403).Render("error", fiber.Map{"msg": "Unauthorised"})
    }
    curruser, err := res.agregator.GetUserById(c.UserContext(), uid)
    if err != nil {
        return c.Status(412).Render("error", fiber.Map{"msg": err})
    }
    users, err := res.agregator.GetUsersWithLimitOffset(c.UserContext(), 1000, 0)
    if err != nil {
        return c.Status(412).Render("error", fiber.Map{"msg": err})
    }
    return c.Render("userlist", fiber.Map{
        "msg": "userlist page: page: Coming soon!",
        "user": curruser,
        "users": users,
    })
}
