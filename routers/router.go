package routers

import (
	"xcfg-server/controllers"
	"xcfg-server/models"

	"github.com/astaxie/beego/context"

	"strings"

	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.MainController{})

	beego.Router("/auth/login", &controllers.AuthController{}, "post:Login")
	beego.Router("/auth/checkpwd", &controllers.AuthController{}, "post:CheckPwd")
	beego.Router("/auth/mod", &controllers.AuthController{}, "post:ModPwd")
	beego.Router("/auth/menus", &controllers.AuthController{}, "get:GetMenus")

	beego.Router("/user/list", &controllers.UserController{}, "post:GetUsers")
	beego.Router("/user/disable", &controllers.UserController{}, "post:DisableUser")
	beego.Router("/user/enable", &controllers.UserController{}, "post:EnableUser")
	beego.Router("/user/add", &controllers.UserController{}, "post:AddUser")
	beego.Router("/user/exists", &controllers.UserController{}, "post:IsExists")
	beego.Router("/user/get", &controllers.UserController{}, "get:GetUser")
	beego.Router("/user/roles", &controllers.UserController{}, "get:GetUserRoles")
	beego.Router("/user/update", &controllers.UserController{}, "post:UpdateUser")

	beego.Router("/role/list", &controllers.RoleController{}, "post:GetRoles")
	beego.Router("/role/add", &controllers.RoleController{}, "post:AddRole")
	beego.Router("/role/exists", &controllers.RoleController{}, "post:IsExists")
	beego.Router("/role/get", &controllers.RoleController{}, "get:GetRole")
	beego.Router("/role/update", &controllers.RoleController{}, "post:UpdateRole")
	beego.Router("/role/resources", &controllers.RoleController{}, "get:GetRoleMenus")

	beego.Router("/menu/list", &controllers.RoleController{}, "post:GetMenus")

	beego.Router("/ConfigVersionHandler.ashx", &controllers.CfgController{}, "post:Post")
	beego.Router("/cfg/list", &controllers.CfgController{}, "post:GetCfgs")
	beego.Router("/cfg/add", &controllers.CfgController{}, "post:AddCfg")
	beego.Router("/cfg/exists", &controllers.CfgController{}, "post:Exists")
	beego.Router("/cfg/get", &controllers.CfgController{}, "get:GetCfg")
	beego.Router("/xcfg/get", &controllers.CfgController{}, "get:Get")
	beego.Router("/cfg/update", &controllers.CfgController{}, "post:UpdateCfg")

	var FilterUser = func(ctx *context.Context) {
		if ctx.Request.RequestURI == "/notfound" || ctx.Request.RequestURI == "/" {
			return
		}
		isPass := false
		//白名单
		menus := models.GetWhiteListByMenuIds()
		for _, v := range menus {
			//fmt.Println(v.Url, ctx.Request.RequestURI)
			if strings.HasPrefix(ctx.Request.RequestURI, v.Url) {
				isPass = true
			}
		}

		if isPass {
			return
		}
		//需要登录
		adAuthCookie := ctx.Input.Cookie("adAuthCookie")
		if adAuthCookie != "true" && ctx.Request.RequestURI != "/" {
			ctx.Redirect(302, "/")
		}
		//用户权限
		un := ctx.GetCookie("loginUser")
		if un != "admin" && ctx.Input.Session("userId") != nil {

			id := ctx.Input.Session("userId").(int)
			roles := models.GetRolesByUserId(id)
			menuIds := []int{0}
			for _, v := range roles {
				rms := models.GetMenusByRoleId(v.RoleId)
				for _, rm := range rms {
					menuIds = append(menuIds, rm.MenuId)
				}
			}
			menus := models.GetResourcesByMenuIds(menuIds)
			for _, v := range menus {
				if strings.HasPrefix(ctx.Request.RequestURI, v.Url) {
					isPass = true
				}
			}
			if !isPass {
				ctx.Redirect(302, "/notfound")
			}
		}

	}

	beego.InsertFilter("/*", beego.BeforeRouter, FilterUser)
}
