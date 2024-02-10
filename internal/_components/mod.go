package _components

import "fmt"

func Banner(hasIcon bool) string {
	bars := ""
	if hasIcon {
		bars = SvgIcon("/static/svg/bars-dark.svg", "sm")
	}
	return fmt.Sprintf(`
		<header class='flex flex-row justify-between items-center p-6'>
			<h2 class='text-xl'>Repairs Log</h2>
			%s
		</header>
	`, bars)
}

func SvgIcon(src string, size string) string {
	var svgSize string
	if size == "sm" {
		svgSize = "h-6 w-6"
	}
	if size == "md" {
		svgSize = "h-8 w-8"
	}
	return fmt.Sprintf(`
		<div class='%s flex items-center cursor-pointer'>
			<img src='%s'/>
		</div>
	`, svgSize, src)
}

func MqWrapperCentered(content string) string {
	return fmt.Sprintf(`
		<div class='flex justify-center items-center w-full h-full'>
			%s
		</div>
	`, content)
}

func FormInputLabel(name string, serverName string) string {
	return fmt.Sprintf(`
		<div class='flex flex-col gap-2 w-full'>
			<label class="text-xs leading-none">%s</label>
			<input name='%s' type='text' value='' autocomplete='off' class="flex p-1 w-full rounded-md border border-darkgray border-input bg-black px-2 py-1 text-sm ring-offset-background placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2">
		</div>
	`, name, serverName)
}

func FormSubmitButton() string {
	return `
		<input type='submit' class="flex w-full row-start-5 row-end-5 col-start-1 col-end-8 bg-white text-black mt-4 items-center justify-center whitespace-nowrap rounded-md text-xs font-medium ring-offset-background transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 bg-primary text-primary-foreground hover:bg-primary/90 h-10 px-4 py-2" type="submit"/ >
	`
}

func FormError(err string) string {
	return fmt.Sprintf(`
		<p class='text-red font-sm'>%s</p>
	`, err)
}

func FormTitle(title string) string {
	return fmt.Sprintf(`
		<h2 class='text-lg'>%s</h2>
	`, title)

}

func LoginForm(err string) string {
	return fmt.Sprintf(`
		<form class='flex flex-col p-6 gap-4 w-full max-w-[500px]'>
			%s%s%s%s%s
		</form>
	`, FormTitle("Login"), FormError(err), FormInputLabel("Username", "username"), FormInputLabel("Password", "password"), FormSubmitButton())
}

