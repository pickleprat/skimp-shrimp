package _components

import (
	"cfasuite/internal/_model"
	"fmt"
)

func Banner(hasIcon bool, menu string) string {
	bars := ""
	x := ""
	if hasIcon {
		bars = SvgIcon("/static/svg/bars-dark.svg", "sm", "x-show='!navopen' @click='navopen = !navopen'")
		x = SvgIcon("/static/svg/x-dark.svg", "sm", "x-show='navopen' x-cloak @click='navopen = !navopen'")
	}
	return fmt.Sprintf(`
		<header x-data="{ navopen: false }">
			<div class='flex flex-row justify-between items-center p-6 h-[75px] z-40 fixed top-0 w-full border-b border-darkgray'>
				<h2 class='text-xl'>Repairs Log</h2>
				%s%s
			</div>
			%s
		</header>
		<div class='h-[75px]'></div>
	`, bars, x, menu)
}

func NavMenuLink(name string, href string, currentPath string) string {
	activeClass := ""
	if currentPath == href {
		activeClass = "underline"
	}
	return fmt.Sprintf(`
		<li>
			<a href='%s' class='%s'>%s</a>
		</li>
	`, href, activeClass, name)
}

func AppNavMenu(currentPath string) string {
	return fmt.Sprintf(`
		<nav x-show='navopen' x-transition x-cloak class='fixed left-0 top-[75px] h-full w-3/5 max-w-[300px] p-6 z-40 bg-black border-r border-darkgray'>
			<ul class='flex flex-col gap-2'>
				%s%s
			</ul>
		</nav>
		<div x-show='navopen' @click='navopen = !navopen' class='h-full w-full fixed top-0 left-0 bg-black opacity-50 z-30'></div>
	`, NavMenuLink("Home", "/app", currentPath), NavMenuLink("Logout", "/logout", currentPath))
}

func SvgIcon(src string, size string, xattr string) string {
	var svgSize string
	if size == "sm" {
		svgSize = "h-6 w-6"
	}
	if size == "md" {
		svgSize = "h-8 w-8"
	}
	return fmt.Sprintf(`
		<div %s class='%s flex items-center cursor-pointer'>
			<img src='%s'/>
		</div>
	`, xattr, svgSize, src)
}

func MqWrapperCentered(content string) string {
	return fmt.Sprintf(`
		<div class='flex justify-center items-center w-full h-full'>
			%s
		</div>
	`, content)
}

func MqGridTwoColEvenSplit(content1 string, content2 string) string {
	return fmt.Sprintf(`
		<div class='grid sm:grid-cols-2 sm:gap-4'>
			%s%s
		</div>
	`, content1, content2)

}

func FormInputLabel(name string, serverName string, inputType string, value string) string {
	if inputType == "" {
		inputType = "text"
	}
	spanText := ""
	if serverName == "phone" {
		spanText = "(format: 123-456-7890)"
	
	}
	return fmt.Sprintf(`
		<div class='flex flex-col gap-2 w-full'>
			<label class="text-xs leading-none">%s <span class='text-gray'>%s</span></label>
			<input name='%s' type='%s' value='%s' autocomplete='off' class="flex p-1 w-full rounded-md border border-darkgray border-input bg-black px-2 py-1 text-sm ring-offset-background placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2">
		</div>
	`, name, spanText, serverName, inputType, value)
}

func FormSubmitButton() string {
	return `
		<input @click='loading = true' x-show='!loading' type='submit' class="flex w-full row-start-5 row-end-5 col-start-1 col-end-8 bg-white text-black mt-4 items-center justify-center whitespace-nowrap rounded-md text-xs font-medium ring-offset-background transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 bg-primary text-primary-foreground hover:bg-primary/90 h-10 px-4 py-2" type="submit"/ >
	`
}

func FormDeleteButton() string {
	return `
		<input @click='loading = true' x-show='!loading' type='submit' class="flex w-full row-start-5 row-end-5 col-start-1 col-end-8 bg-red text-black mt-4 items-center justify-center whitespace-nowrap rounded-md text-xs font-medium ring-offset-background transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 bg-primary text-white text-primary-foreground hover:bg-primary/90 h-10 px-4 py-2" type="submit" value='Delete'/>
	`
}

func FormError(err string) string {
	if err == "" {
		return ""
	
	} else {
		return fmt.Sprintf(`
			<p class='text-red text-xs py-2'>%s</p>
		`, err)
	}
}

func FormTitle(title string) string {
	return fmt.Sprintf(`
		<h2 class='text-lg'>%s</h2>
	`, title)
}

func FormLoader() string {
	return `
		<div x-show='loading' x-cloak class='flex items-center justify-center mt-4'>
			<div class='h-10 w-10 rounded-full border-2 border-gray border-t-white animate-spin'></div>
		</div>
	`
}

func LoginForm(err string, username string, password string) string {
	return fmt.Sprintf(`
		<form x-data="{ loading: false }" method='POST' class='flex flex-col p-6 gap-4 w-full max-w-[500px]'>
			%s%s%s%s%s%s
		</form>
	`, FormTitle("Login"), FormError(err), FormInputLabel("Username", "username", "", username), FormInputLabel("Password", "password", "password", password), FormLoader(), FormSubmitButton())
}

func CreateManufacturerForm(err string, name string, phone string, email string, xclass string) string {
	return fmt.Sprintf(`
		<form x-data="{ loading: false }" method='POST' class='flex flex-col p-6 gap-4 w-full max-w-[500px] %s'>
			%s%s%s%s%s%s%s
		</form>
	`, xclass, FormTitle("Create Manufacturer"), FormError(err), FormInputLabel("Name", "name", "", name), FormInputLabel("Phone", "phone", "", phone), FormInputLabel("Email", "email", "", email), FormLoader(), FormSubmitButton())
}

func ManufacturerList(manufacturers []_model.Manufacturer, xclass string) string {
	manufacturerList := ""
	for _, manufacturer := range manufacturers {
		href := fmt.Sprintf("/app/manufacturer/%d", manufacturer.ID)
		manufacturerList += fmt.Sprintf(`
			<tr>
				<td>
					<a href='%s' class='underline select-none hover:text-red cursor-pointer'>%s</a>
				</td>
				<td>%s</td>
				<td class='hidden sm:flex'>%s</td>
			</tr>
		`, href, manufacturer.Name, manufacturer.Phone, manufacturer.Email)
	}
	if len(manufacturers) == 0 {
		return fmt.Sprintf(`
			<div class='p-6 max-w-[500px] w-full %s'>
				<p>Oh, so empty!</p>
			</div>
		`, xclass)
	} else {
		return fmt.Sprintf(`
			<div class='p-6 max-w-[500px] w-full %s'>
				<table class='text-xs w-full'>
					<tr class='text-left'>
						<th>Name</th>
						<th>Phone</th>
						<th class='hidden sm:flex'>Email</th>
					</tr>
					%s
				</table>
			</div>
		`, xclass, manufacturerList)
	}
}

func ManufacturerDetails(manufacturer _model.Manufacturer, update string, deleteErr string, updateErr string, xclass string) string {
	if update == "" {
		update = "true"
	
	}
	deleteAction := fmt.Sprintf("/app/manufacturer/%d/delete", manufacturer.ID)
	updateActon := fmt.Sprintf("/app/manufacturer/%d/update", manufacturer.ID)
	return fmt.Sprintf(`
		<div x-data="{ update: %s }" class='p-6 w-full flex flex-col gap-4 max-w-[500px] %s'>
			<div class='flex flex-row justify-between'>
				<h2>Manufacturer Details</h2>
				%s%s
			</div>
			<div class='text-xs'>
				<p><strong>Name:</strong> %s</p>
				<p><strong>Phone:</strong> %s</p>
				<p><strong>Email:</strong> %s</p>
			</div>
			<form x-cloak x-show='!update' method='POST' action='%s' x-data="{ loading: false }">
				<h2 class='mb-4'>Delete Manufacturer</h2>
				%s%s%s%s
			</form>
			<form x-cloak x-show='update' method="POST" action="%s" class='flex flex-col gap-4' x-data="{ loading: false }">
				<h2>Update Manufacturer</h2>
				%s%s%s%s%s%s
			</form>
		</div>
	`, 
		update,
		xclass, 
		SvgIcon("/static/svg/edit-dark.svg", "sm", "x-cloak x-show='!update ' @click='update=!update'"), 
		SvgIcon("/static/svg/delete-dark.svg", "sm", "x-cloak x-show='update' @click='update=!update'"), 
		manufacturer.Name, 
		manufacturer.Phone, 
		manufacturer.Email, 
		deleteAction, 
		FormError(deleteErr), 
		FormInputLabel("Type '" + manufacturer.Name + "' to delete", "name", "", ""), 
		FormDeleteButton(), 
		FormLoader(), 
		updateActon,
		FormError(updateErr),
		FormInputLabel("Name", "name", "", manufacturer.Name),
		FormInputLabel("Phone", "phone", "", manufacturer.Phone), 
		FormInputLabel("Email", "email", "", manufacturer.Email),
		FormLoader(),
		FormSubmitButton(),
	)
}

func CreateEquipmentForm(err string, xclass string) string {
	return fmt.Sprintf(`
		<form x-data="{ loading: false }" method='POST' class='flex flex-col p-6 gap-4 w-full max-w-[500px] %s'>
			%s%s%s%s
			<div class='flex flex-col text-xs w-fit rounded gap-2'>
				<label>Photo</label>
				<button type='button' x-on:click="$refs.uploadInput.click(); console.log('yo')" class='text-left border border-white p-2 rounded'>Upload Photo</button>
				<input x-ref='uploadInput' x-on:change='console.log("changed")' type='file' class='hidden'/>
			</div>
			%s%s
		</form>
	`, 
		xclass, 
		FormTitle("Create Equipment"), 
		FormError(err), 
		FormInputLabel("Nickname", "nickname", "", ""), 
		FormInputLabel("Serial Number", "number", "", ""), 
		FormLoader(),
		FormSubmitButton(),
	)
}

