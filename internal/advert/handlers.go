package advert

import (
    "strings"
    "errors"
    "time"
    // "fmt"
    recapcheck "github.com/dpapathanasiou/go-recaptcha"
    "github.com/gofiber/fiber/v2"
    "github.com/tvitcom/local-adverts/internal/config"
    "github.com/tvitcom/local-adverts/pkg/util"  
)

func (res resource) pageIndex(c *fiber.Ctx) error {
    var adverts []Advert
    searchclause := strings.TrimSpace(c.Query("q", ""))
    if searchclause != "" {
        form := new(QuickSearchForm) 
        form.Clause = searchclause
        // Form validation    
        if err := form.Validate(); err != nil {
            return c.Status(412).Render("error", fiber.Map{"msg": err})
        }
        res, err := res.agregator.GetAdvertsSearch(c.UserContext(), form.Clause)
        if err != nil {
            return c.Status(500).Redirect("/error.html?msg=Ошибка работы сайта")
        }
        adverts = res
    } else {
        res, err := res.agregator.GetAdvertsLast(c.UserContext())
        if err != nil {
            return c.Status(500).Redirect("/error.html?msg=Ошибка работы сайта")
        }
        adverts = res
    }
    return c.Render("index", fiber.Map{
        "msg": "index page",
        "adverts": adverts,
    })
}

func (res resource) pageWatch(c *fiber.Ctx) error {
    uid := util.Pkeyer(c.Locals("iam"))
    advid := util.Pkeyer(c.Query("advid", "0"))
    if advid == 0 {
        res.logger.With(c.UserContext()).Error("Ошибка GET advid параметра")
        return c.Status(500).Redirect("/error.html?msg=Ошибка параметра запроса.")
    }
    advert, err := res.agregator.GetAdvertById(c.UserContext(), advid)
    if err != nil {
    println(err.Error())
        return c.Status(404).Redirect("/error.html?msg=Объявление не найдено или удалено.")
    }
    author, err := res.agregator.GetUserById(c.UserContext(), advert.UserId)
    if err != nil {
        return c.Status(500).Redirect("/error.html?msg=Ошибка работы сайта.")
    }
    return c.Render("watch", fiber.Map{
        "msg": "watch page",
        "recaptcha_site_key": config.CFG.RecaptchaSiteKey,
        "advert": advert,
        "user": author,
        "fqdn": config.CFG.AppFqdn,
        "uid": uid,
    })
}

func (res resource) handlerWatchAuthor(c *fiber.Ctx) error {
    form := new(WatchAuthorForm)
    form.SignerUA = c.Cookies("privacy_signer_ua", "")
    form.SignerScreen = c.Cookies("privacy_signer_screen", "")
    form.SignerLangs = c.Cookies("privacy_signer_langs", "")
    form.SignerTime = c.Cookies("privacy_signer_time", "")
    if err := c.BodyParser(form); err != nil {
        res.logger.With(c.UserContext()).Error(err.Error())
        return c.Status(412).JSON(&fiber.Map{
            "ok": false,
            "data": "Переданы невалидные данные",
        })
    }
// println("FORM_RECAPTCHA_RESPONSE:", form.RecaptchaResponse)
    // Validation - Stage 1 of 2
    if form.SignerUA == "" || form.SignerScreen == "" || form.SignerLangs == "" || form.SignerTime == "" {
        // return c.Status(403).Render("error",fiber.Map{"msg": "Пожалуйста ознакомьтесь с информацией о конфиденциальности сайта и примите решение об использовании нашего сайта!"})
        return c.Status(403).JSON(&fiber.Map{
            "ok": false,
            "data": "Пожалуйста ознакомьтесь с информацией о конфиденциальности сайта и примите решение об использовании нашего сайта!",
        })
    }
// println("FORM:",form.AdvertId, form.SignerUA, form.SignerScreen, form.SignerLangs, form.SignerTime)
    // Validation - Stage 2 of 2
    if err := form.Validate(); err != nil {
        var errMsg string
        if config.CFG.AppMode == "dev" {
            errMsg = err.Error()
        } else {
            errMsg = "Невалидный запрос"
        }
        return c.Status(403).JSON(&fiber.Map{
            "ok": false,
            "data": errMsg,
        })
    }
    var result bool 
    //!!!  if config.CFG.AppMode != "dev" {
    if false {
        // Make google recaptcha-validation service request
        recapcheck.Init(config.CFG.RecaptchaSecret)
        res, err := recapcheck.Confirm(c.IP(), form.RecaptchaResponse)
        if err != nil {
            return c.Status(500).JSON(&fiber.Map{
                "ok": false,
                "data": err.Error(),
            })
        }
        result = res
    } else {
        result = true //!!! Fake specially skip requests to Google for dev mode
    }
    if result {
        author, err := res.agregator.GetUserByAdvertId(c.UserContext(), form.AdvertId)
        if err != nil {
            return c.Status(500).JSON(&fiber.Map{
                "ok": result,
                "data": "Ошибка работы сайта",
            })
        }
        return c.JSON(&fiber.Map{
            "ok": result,
            "data": author.Tel,
        })
    }
    return c.Status(403).JSON(&fiber.Map{
        "ok": result,
        "data": "не доступно",
    })
}

func (res resource) handlerActivity(c *fiber.Ctx) error {
    uid := util.Pkeyer(c.Query("uid", "0"))
    if uid == 0 {
        res.logger.With(c.UserContext()).Error("Ошибка GET advid параметра")
        return c.Status(500).Redirect("/error.html?msg=Ошибка параметра")
    }
    adverts, err := res.agregator.GetAdvertsByUserId(c.UserContext(), uid)
    if err != nil {
        return c.Status(500).Redirect("/error.html?msg=Ошибка работы сайта")
    }
    author, err := res.agregator.GetUserById(c.UserContext(), uid)
    if err != nil {
        return c.Status(500).Redirect("/error.html?msg=Ошибка работы сайта")
    }
    return c.Render("activity", fiber.Map{
        "msg": "activity page",
        "adverts": adverts,
        "author": author,
    })
}

func (res resource) pagePublication(c *fiber.Ctx) error {
    uid := util.Pkeyer(c.Locals("iam"))
    user, err := res.agregator.GetUserById(c.UserContext(), uid)
    if err != nil {
        return c.Status(500).Redirect("/error.html?msg=Ошибка работы сайта")
    }
    categories, err := res.agregator.GetCategoriesPath(c.UserContext())
    if err != nil {
        return c.Status(500).Redirect("/error.html?msg=Ошибка работы сайта")
    }
    return c.Render("publication", fiber.Map{
        "msg": "avdert page",
        "user": user,
        "categories": categories,
    })
}

// ->makeUser->makeAdvertInactive->sendApproveUser1 [->handleEmailApprove]->userprofile->useradverts
func (res resource) handlerPublication(c *fiber.Ctx) error {
    // format in the assets/media/[advert_id]_[num].jpg pictures format
    handle_dt := time.Now().Format("2006-01-02 15:04:05") // Datetime for seve resources by bigint(datetime)
    uid := util.Pkeyer(c.Locals("iam"))
    logined := uid
    form := new(QuickAdvertForm) 
    if err := c.BodyParser(form); err != nil {
        res.logger.With(c.UserContext()).Error(err.Error())
        return c.Status(412).Render("error", fiber.Map{"msg": err})
    }
    // Form validation    
    if err := form.Validate(); err != nil {
        return c.Status(412).Render("error", fiber.Map{"msg": err})
    }
    // Get the UserId from db or make new:
    var user User
    if uid == 0 {
        // Find user by email
        user, _ = res.agregator.GetUserByEmail(c.UserContext(), form.Email)
        if user.UserId == 0 {
            // Save New User record
            userForm := &SignupForm{
                Email:     form.Email,
                GivenName: form.GivenName,
                NewPassword:       handle_dt,
                NewPasswordRepeat: handle_dt,
                Tel: form.Tel,
            }
            var err error
            user.UserId, err = res.agregator.CreateUser(c.UserContext(), userForm, handle_dt, "user", "byquickadv", "")
            if err != nil {
                res.logger.With(c.UserContext()).Error(err.Error())
                return c.Status(412).Render("error", fiber.Map{"msg": "CreateUser err"})
            }
        }
        // Найти текущего пользователя и проверить соответствие
        findeduser, errU := res.agregator.GetUserByEmail(c.UserContext(), form.Email)
        if errU != nil {
            return c.Status(403).Render("error", fiber.Map{"msg": "errU"})
        }
        user = findeduser
        uid = user.UserId
    }

    // Save new Category record
    // Save new Advert inactive record
    advertId, err := res.agregator.CreateAdvert(c.UserContext(), form, uid, handle_dt)
    if err != nil {
        res.logger.With(c.UserContext()).Info(err.Error())
        return c.Status(500).Render("error", fiber.Map{"msg": "Ошибка записи объявления в БД"})
    }
    // Parse the multipart form:
    if ff, err := c.MultipartForm(); err == nil {

        picformfields := []string{"1", "2", "3", "4"} //free service - pictures
        for _, v := range picformfields {

            files := ff.File["picture" + v]// => []*multipart.FileHeader
            // Loop through files:
            for i, ff := range files {
                if i > 1 {
                    return errors.New("Файлов слишком много для сайта")
                }

                // Start Image convey:      
                imagerawfname := util.Stringer(advertId) + "_" + v + "_raw.jpg"
                imagefname := util.Stringer(advertId) + "_" + v +".jpg"
                
                if err := c.SaveFile(ff, config.PictureAdvertsPath + imagerawfname); err != nil {
                    return c.Status(500).Render("error", fiber.Map{
                        "msg": err,
                    })
                }
                if err := util.ImagefileValidations(config.PictureAdvertsPath + imagerawfname); err != nil {
                    return c.Status(501).Render("error", fiber.Map{
                        "msg": "Загруженная картинка не подходит для сайта",
                    })
                }
                if err := util.ImagefileResizing(config.PictureAdvertsPath + imagerawfname, config.PictureAdvertsPath + imagefname, 468); err != nil {
                    return c.Status(500).Render("error", fiber.Map{
                        "msg": "Загруженная картинка не подходит для обработки.",
                    })
                }
                err = util.ImagefileProgressiveOptimisation(c.UserContext(), config.PictureAdvertsPath + imagefname, "", true)
                if err != nil {
                    return c.Status(500).Render("error", fiber.Map{
                        "msg": "Загруженная картинка не обработана.",
                    })
                }
                // Update Advert record with pictures names
                err = res.agregator.UpdateAdvertsPicture(c.UserContext(), advertId, "picture" + v, imagefname)
                if err != nil {
                    return c.Status(500).Render("error", fiber.Map{
                        "msg": "Загруженная картинка для пользователя не сохранена.",
                    })
                }
            }
        }
    }
    /* LOGIC: 
        user: msg "Thanks" , exit.

        guest: if old-user: msg: "Thanks", exit
               else: mail-activation-acc(), msg: "Activation link sended to you email", exit 
    */
    var msg string
    if logined > 0 {
        msg = "Благодарим! Теперь с помощью электронного кабинета сможете актуализировать в поиске и удалять свои объявления"
    } else {
        msg = "Теперь вы сможете публиковать объявления использую свой указанный e-mail"
        // msg = "Подтвердите объявление в течении часа по отправленной ссылке на ваш e-mail"
        // var tplParams map[string]string
        // tplParams["From"] = config.CFG.MailSmtphost
        // tplParams["Name"] = user.Name
        // tplParams["Brand"] = config.CFG.BizName
        // tplParams["ApproveURL"] = "https://" + config.CFG.AppFqdn + "/100500"
        // // Notification by Email (CTX, smtpHost, smtpPort, fromMail, fromPassword, toMail, tplFile string, mailData map[string]string) error 
        // if err := util.SendEmailWithTempate(
        //     c.UserContext(), 
        //     config.CFG.MailSmtphost, 
        //     config.CFG.MailSmtpport,
        //     config.CFG.MailUsername, 
        //     config.CFG.MailPassword, 
        //     user.Email, 
        //     "web/email/approvement.emlt", 
        //     tplParams,
        // ); err != nil {
        //     return c.Status(500).Render("error", fiber.Map{"msg": "Ошибка отправки подтверждения отправки. Обратитесь в поддержку сайта"})
        // }
    }
    return c.Status(201).Render("thanks", fiber.Map{
        "msg": msg,
    })
}

func (res resource) pageError(c *fiber.Ctx) error {  
    return c.Render("error", fiber.Map{
        "msg": c.Query("msg", "Чтото пошло не так."),
    })
}

func (res resource) pageMessage(c *fiber.Ctx) error {
    uid := util.Pkeyer(c.Locals("iam"))
    user, err := res.agregator.GetUserById(c.UserContext(), uid)
    if err != nil {
        return c.Status(500).Redirect("/error.html?msg=Ошибка работы сайта.")
    }
    aid := util.Pkeyer(c.Query("advid", "0"))
    advert, err := res.agregator.GetAdvertById(c.UserContext(), aid)
    if err != nil {
        return c.Status(500).Redirect("/error.html?msg=Ошибка работы сайта.")
    }
    advertsuser, err := res.agregator.GetUserById(c.UserContext(), advert.UserId)
    if err != nil {
        return c.Status(500).Redirect("/error.html?msg=Ошибка работы сайта.")
    }
    return c.Render("message", fiber.Map{
        "msg": "Сообщение пользователю",
        "advertsuser": advertsuser,
        "user": user,
        "aid": aid,
    })
}

func (res resource) handlerMessage(c *fiber.Ctx) error {
    // format in the [userId]/[datetime]_[123].jpg pictures format
    handle_dt := time.Now().Format("2006-01-02 15:04:05") // Datetime for seve resources by bigint(datetime)
    uid := util.Pkeyer(c.Locals("iam"))
    
    form := new(MessageForm) 
    if err := c.BodyParser(form); err != nil {
        res.logger.With(c.UserContext()).Info(err)
        return c.Status(412).Render("error", fiber.Map{"msg": err})
    }
    if err := form.Validate(); err != nil {
        return c.Status(412).Render("error", fiber.Map{"msg": err})
    }

    // Get the UserId from db or make new:
    var user User
    if uid == 0 {
        // Find user by email
        user, _ = res.agregator.GetUserByEmail(c.UserContext(), form.Email)
        if user.UserId == 0 {
            // Save New User record
            userForm := &SignupForm{
                Email:     form.Email,
                GivenName: form.GivenName,
                NewPassword:       handle_dt,
                NewPasswordRepeat: handle_dt,
                Tel: "",
            }
            _, err := res.agregator.CreateUser(c.UserContext(), userForm, handle_dt, "user", "regbysupport", "")
            if err != nil {
                res.logger.With(c.UserContext()).Info(err)
                return c.Status(412).Render("error", fiber.Map{"msg": "CreateUser err"})
            }
        }
        // Найти текущего пользователя и проверить соответствие
        findeduser, errU := res.agregator.GetUserByEmail(c.UserContext(), form.Email)
        if errU != nil || form.Email != findeduser.Email || form.GivenName != findeduser.Name {
            return c.Status(403).Render("error", fiber.Map{"msg": "errU"})
        }
        user = findeduser
    }
    advert, err := res.agregator.GetAdvertById(c.UserContext(), form.AdvertId)
    if err != nil {
        return c.Status(500).Redirect("/error.html?msg=Ошибка работы сайта.")
    }
    if err := res.agregator.CreateMessage(c.UserContext(), user.UserId, advert.UserId, form.Msg, handle_dt); err != nil {
        return c.Status(500).Render("error", fiber.Map{"msg": "Ошибка создания сообщения."})
    }
    return c.Status(201).Render("thanks", fiber.Map{
        "msg": "Возможно автор объявления ответит вам в ближайшее время.",
    })
}

func (res resource) pageSupport(c *fiber.Ctx) error {
    uid := util.Pkeyer(c.Locals("iam"))
    aid := c.Query("advid", "0")
    subject := c.Query("subj","tech")
    user, err := res.agregator.GetUserById(c.UserContext(), uid)
    if err != nil {
        return c.Status(500).Redirect("/error.html?msg=Ошибка работы сайта.")
    }
    return c.Render("support", fiber.Map{
        "msg": "support page: Coming soon",
        "subject": subject,
        "user": user,
        "advid": aid,
    })
}

func (res resource) handlerSupport(c *fiber.Ctx) error {
    // format in the [userId]/[datetime]_[123].jpg pictures format
    handle_dt := time.Now().Format("2006-01-02 15:04:05") // Datetime for seve resources by bigint(datetime)
    uid := util.Pkeyer(c.Locals("iam"))
    
    form := new(SupportForm) 
    if err := c.BodyParser(form); err != nil {
        res.logger.With(c.UserContext()).Info(err)
        return c.Status(412).Render("error", fiber.Map{"msg": err})
    }

    // Form validation    
    if err := form.Validate(); err != nil {
        return c.Status(412).Render("error", fiber.Map{"msg": err})
    }

    // Get the UserId from db or make new:
    var user User
    if uid == 0 {
        // Find user by email
        user, _ = res.agregator.GetUserByEmail(c.UserContext(), form.Email)
        if user.UserId == 0 {
            // Save New User record
            userForm := &SignupForm{
                Email:     form.Email,
                GivenName: form.GivenName,
                NewPassword:       handle_dt,
                NewPasswordRepeat: handle_dt,
                Tel: "",
            }
            _, err := res.agregator.CreateUser(c.UserContext(), userForm, handle_dt, "user", "regbysupport", "")
            if err != nil {
                res.logger.With(c.UserContext()).Info(err)
                return c.Status(500).Render("error", fiber.Map{"msg": "CreateUser err"})
            }
        }
        // Найти текущего пользователя и проверить соответствие
        findeduser, errU := res.agregator.GetUserByEmail(c.UserContext(), form.Email)
        if errU != nil {
            return c.Status(500).Render("error", fiber.Map{"msg": "Ошибка работы сайта"})
        }
        user = findeduser
    } else {
        // Найти текущего пользователя по iam/uid
        findeduser, errU := res.agregator.GetUserById(c.UserContext(), uid)
        if errU != nil {
            return c.Status(403).Render("error", fiber.Map{"msg": "errU"})
        }
        user = findeduser
    }

    // Save new Users Message record
    // Parse the multipart form:
    imagepartialfname := util.PhoneNormalisation(handle_dt)
    if ff, err := c.MultipartForm(); err == nil {

        picformfields := [3]string{"1", "2", "3"}
        for _, v := range picformfields {

            files := ff.File["picture" + v]// => []*multipart.FileHeader

            // Loop through files:
            for i, ff := range files {
                if i > 1 {
                    return errors.New("Файлов слишком много послано для сайта.")
                }
                // Start Image convey:      
                userpicfpath, err := util.MakeUploadDirByUserId(config.PictureSupportPath, util.Stringer(user.UserId))
                if err != nil {
                    return c.Status(501).Render("error", fiber.Map{
                        "msg": err,
                    })
                }
                imagerawfname := imagepartialfname + "_raw_" + v +".jpg"
                imagefname := imagepartialfname + "_" + v +".jpg"
                
                if err := c.SaveFile(ff, userpicfpath + imagerawfname); err != nil {
                    return c.Status(500).Render("error", fiber.Map{
                        "msg": err,
                    })
                }

                if err := util.ImagefileValidations(userpicfpath + imagerawfname); err != nil {
                    return c.Status(501).Render("error", fiber.Map{
                        "msg": "Загруженная картинка не подходит для сайта.",
                    })
                }
                err = util.ImagefileProgressiveOptimisation(c.UserContext(), userpicfpath + imagerawfname, userpicfpath + imagefname, false)
                if err != nil {
                    return c.Status(500).Render("error", fiber.Map{
                        "msg": "Загруженная картинка не обработана.",
                    })
                }
            }
        }
    }
    if err := res.agregator.CreateMessage(c.UserContext(), user.UserId, config.SupportUserID, "[" + form.Subject + "]:" + form.Msg, handle_dt); err != nil {
        return c.Status(500).Render("error", fiber.Map{"msg": "Ошибка создания записи запроса."})
    }
    return c.Status(201).Render("thanks", fiber.Map{
        "msg": "Ответим по вашему запросу на ваш <" + form.Email + "> в ближайшее время.",
    })
}

func (res resource) pageThanks(c *fiber.Ctx) error {
    return c.Render("thanks", fiber.Map{
            "thanks": c.Query("msg", "Благодарим за использование нашего сайта " + config.CFG.BizName),
        })
}

func (res resource) pageSoon(c *fiber.Ctx) error {
    return c.Render("soon", fiber.Map{
        "msg": "Контент станет доступен позже.",
    })
}

func (res resource) pageAgreement(c *fiber.Ctx) error {
    return c.SendFile("./assets/docs/agreement.txt", true)
}

func (res resource) pageGdprPolicy(c *fiber.Ctx) error {
    return c.SendFile("./assets/docs/GDPR_POLICY_RU.txt", true)
}

func (res resource) pageRobots(c *fiber.Ctx) error {
    return c.SendFile(config.RobotsFilePath, true)
}

func (res resource) pageSitemap(c *fiber.Ctx) error {
    return c.SendFile(config.SitemapFilePath, true)
}

func (res resource) pagePayment(c *fiber.Ctx) error {  
    return c.Render("payment", fiber.Map{
        "msg": c.Query("msg", "Страница на которой можно оплатить пока не создана. Извините."),
    })
}

func (res resource) handlerCspCollector(c *fiber.Ctx) error {
    res.logger.With(c.UserContext()).Info(string(c.Body()))
    return c.JSON(&fiber.Map{
            "ok": true,
            "data": "log ok",
        })
}
