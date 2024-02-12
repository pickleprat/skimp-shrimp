package _components

import "fmt"

func Banner(hasIcon bool, menu string) string {
	bars := ""
	x := ""
	if hasIcon {
		bars = SvgIcon("/static/svg/bars-dark.svg", "sm", "x-show='!navopen' @click='navopen = !navopen'")
		x = SvgIcon("/static/svg/x-dark.svg", "sm", "x-show='navopen' x-cloak @click='navopen = !navopen'")
	}
	return fmt.Sprintf(`
		<header x-data="{ navopen: false }" class='flex flex-row justify-between items-center p-6'>
			<h2 class='text-xl'>Repairs Log</h2>
			%s
			%s
		</header>
	`, bars, x)
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

func FormInputLabel(name string, serverName string, inputType string) string {
	if inputType == "" {
		inputType = "text"
	}
	return fmt.Sprintf(`
		<div class='flex flex-col gap-2 w-full'>
			<label class="text-xs leading-none">%s</label>
			<input name='%s' type='%s' value='' autocomplete='off' class="flex p-1 w-full rounded-md border border-darkgray border-input bg-black px-2 py-1 text-sm ring-offset-background placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2">
		</div>
	`, name, serverName, inputType)
}

func FormSubmitButton() string {
	return `
		<input @click='loading = true' x-show='!loading' type='submit' class="flex w-full row-start-5 row-end-5 col-start-1 col-end-8 bg-white text-black mt-4 items-center justify-center whitespace-nowrap rounded-md text-xs font-medium ring-offset-background transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 bg-primary text-primary-foreground hover:bg-primary/90 h-10 px-4 py-2" type="submit"/ >
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

func LoginForm(err string) string {
	return fmt.Sprintf(`
		<form x-data="{ loading: false }" method='POST' class='flex flex-col p-6 gap-4 w-full max-w-[500px]'>
			%s%s%s%s%s%s
		</form>
	`, FormTitle("Login"), FormError(err), FormInputLabel("Username", "username", ""), FormInputLabel("Password", "password", "password"), FormLoader(), FormSubmitButton())
}

