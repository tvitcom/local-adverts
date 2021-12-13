package advert

import (
    "fmt"
    "time"
    "context"
    "os/exec"
    // "errors"
    
    "github.com/nfnt/resize"
    "image/jpeg"
    "os"
    "strings"

    "github.com/gofiber/fiber/v2"
    // "github.com/tvitcom/local-adverts/pkg/log"
    "github.com/tvitcom/local-adverts/internal/config"
    "github.com/tvitcom/local-adverts/pkg/util"  
    "github.com/koyachi/go-nude"
)


func (res resource) pageLogin(c *fiber.Ctx) error {
    return c.Render("login", fiber.Map{
        "msg": "avdert page: Coming soon",
    })
}

func (res resource) handlerLogin(c *fiber.Ctx) error {
    handle_dt := time.Now().Format("2006-01-02 15:04:05")
    form := new(LoginForm)
    if err := c.BodyParser(form); err != nil {
        res.logger.With(c.UserContext()).Error(err.Error())
        return c.Status(412).Render("error", fiber.Map{"msg": err})
    }
    // validation    
    if err := form.Validate(); err != nil {
        return c.Status(403).Render("error", fiber.Map{
            "msg": "Неверно введены данные формы логина",
        })
    }
    
    user, err := res.agregator.GetUserByEmail(c.UserContext(), form.Username)
 
    // currentHash := util.MakeBCryptHash(form.CurrentPassword, config.BCRYPT_COST)
    // fmt.Println("THAT PASSWORD:", currentHash)
    // fmt.Printf("SQL: UPDATE user SET passhash = '%s' WHERE email = '%s';\n", currentHash, form.Username)
    
    if err != nil || util.VerifyBCrypt(form.CurrentPassword, user.Passhash) != nil {
        return c.Status(403).Render("error", fiber.Map{
            "msg": "Неверно ввели логин или пароль или всё вместе",
        })
    }
    // Update the user.lastlogin
    if err := res.agregator.UpdateUserLastlogin(c.UserContext(), user.UserId, handle_dt); err != nil {
        return c.Status(500).Redirect("/error.html?msg=Ошибка работы сайта")
    }
    
    // Let identity marker - user authenticated successfully
    // seckey, tokid, fqdn, uid, appsid, roles
    rnd32 := util.RandomHexString(8)
    tok, err := util.MakeJwtString(config.CFG.AppSecretKey, rnd32, config.CFG.AppFqdn, util.Stringer(user.UserId), "main", "user")
    if err != nil {
        return c.Status(403).Render("error", fiber.Map{"msg": err})
    }
    makeJWTCookie(c, tok)
    return c.Redirect("/my/useradverts.html", 301)
}
func (res resource) handlerLogout(c *fiber.Ctx) error {
    deleteJWTCookie(c)
    return c.Render("thanks", fiber.Map{
        "msg": "за посещение сайта. Удачных сделок!",
    })
    return c.Redirect("/thanks.html", 301)
}
func (res resource) pageSignup(c *fiber.Ctx) error {
    return c.Render("signup", fiber.Map{
        "msg": "Добро пожаловать",
    })
}

func (res resource) handlerSignup(c *fiber.Ctx) error {
    handle_dt := time.Now().Format("2006-01-02 15:04:05") // Datetime for save additional resources by this key 
    form := new(SignupForm)
    if err := c.BodyParser(form); err != nil {
        res.logger.With(c.UserContext()).Error(err.Error())
        return c.Status(412).Render("error", fiber.Map{"msg": err})
    }

    fmt.Println("FORM:", form, fmt.Sprintf("%s", time.Now().UnixNano()))
     
    // validation    
    if err := form.Validate(); err != nil {
        return c.Status(412).Render("error", fiber.Map{"msg": err})
    }

    userfilename := ""
    pictureFile, err := c.FormFile("picture")
    if err != nil {
        return c.Status(412).Render("error", fiber.Map{
            "msg": err,
        })
    }

    if pictureFile != nil {

        // Save file to root directory:fmt.Sprintf("./%s", file.Filename
        fname := fmt.Sprintf(config.UploadedPath + "%s", pictureFile.Filename)
        fileErr := c.SaveFile(pictureFile, fname)
        if fileErr != nil {
            return c.Status(500).Render("error", fiber.Map{
                "msg": fileErr,
            })
        }

        // Erotic photography validation
        isNude, fileErr := nude.IsNude(fname)
        if fileErr != nil {
            return c.Status(500).Render("error", fiber.Map{
                "msg": fileErr,
            })
        } else if isNude {
            return c.Status(500).Render("error", fiber.Map{
                "msg": "Загруженная картинка не подходит для сайта",
            })
        }

        //Resizing to width 65(size for avatars) using Lanczos resampling
        file, err := os.Open(fname)
        if err != nil {
            return c.Status(500).Render("error", fiber.Map{
                "msg": err,
            })
        }
        img, err := jpeg.Decode(file)
        if err != nil {
            return c.Status(500).Render("error", fiber.Map{
                "msg": err,
            })
        }
        file.Close()
        m := resize.Resize(75, 0, img, resize.Lanczos3)
        userfilename = util.GetMD5Hash(fmt.Sprintf("%s", time.Now().UnixNano()))+".jpg"
        resizedfile := config.PictureUserPath + userfilename
        out, err := os.Create(resizedfile)
        if err != nil {
            return c.Status(500).Render("error", fiber.Map{
                "msg": err,
            })
        }
        defer out.Close()
        jpeg.Encode(out, m, nil)

        // Jpegoptimizing
        lsCmd := exec.Command("bash", "-c", "file " + resizedfile)
        lsOut, err := lsCmd.Output()
        if err != nil {
            panic(err)
        }
        if !strings.Contains(string(lsOut), "progressive") {
            ConvertingTimeout := 5 * time.Second
            ctx, cancel := context.WithTimeout(c.UserContext(), ConvertingTimeout)
            defer cancel()
            if err := exec.CommandContext(ctx, "jpegoptim", "--strip-all", "--all-progressive", "-ptm85", "--path="+config.PictureUserPath, resizedfile).Run(); err != nil {
                return c.Status(500).Render("error", fiber.Map{
                    "msg": "CAUSE: ConvertingTimeout occured",
                })
            }
        }
    }

    // Save User record
    _, err = res.agregator.CreateUser(c.UserContext(), form, handle_dt, "user", "regbyform", userfilename)
    if err != nil {
       return err
    }
    /* LOGIC ACTIVATE RECORD:
    if guest-->newuser: send-activate-email(), exit
    */
    return c.Render("login", fiber.Map{
        "msg": "Сейчас попробуйте зайти на сайт с вашим логином и паролем",
    })
}

func (res resource) pageGoogleuser(c *fiber.Ctx) error {
    return c.Render("googleuser", fiber.Map{
        "msg": "Добро пожаловать",
        "gclientid": config.CFG.GoogleClientID,
    })
}

func (res resource) handlerGoogleuser(c *fiber.Ctx) error {
println("RETURNED:")
println(c.Body())
    return c.Render("googleuser", fiber.Map{
        "msg": "Добро пожаловать",
    })
}

func makeJWTCookie(c *fiber.Ctx, jwt string) {
    duration := 120 * time.Minute
    cookie := new(fiber.Cookie)
    cookie.Name = "tok"
    cookie.Value = jwt
    cookie.HTTPOnly = true
    cookie.Expires = time.Now().Add(duration)
    c.Cookie(cookie)
}

func deleteJWTCookie(c *fiber.Ctx) {
    negateDuration := -1 * time.Minute
    cookie := new(fiber.Cookie)
    cookie.Name = "tok"
    cookie.Value = "logout"
    cookie.HTTPOnly = true
    cookie.Expires = time.Now().Add(negateDuration)
    c.Cookie(cookie)
}