package _components

import (
	"cfasuite/internal/_model"
	"cfasuite/internal/_util"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"gorm.io/gorm"
)

func Banner(hasIcon bool, menu string) string {
	bars := ""
	x := ""
	if hasIcon {
		bars = SvgIcon("/static/svg/bars-dark.svg", "sm", "x-show='!navopen' @click='navopen = !navopen'", "")
		x = SvgIcon("/static/svg/x-dark.svg", "sm", "x-show='navopen' x-cloak @click='navopen = !navopen'", "")
	}
	return fmt.Sprintf(`
		<header x-data="{ navopen: false }">
			<div class='flex flex-row justify-between items-center p-6 h-[75px] z-40 bg-black fixed top-0 w-full border-b border-darkgray'>
				<h2 class='text-xl'>Repairs Log</h2>
				%s%s
			</div>
			%s
		</header>
		<div class='h-[75px]'></div>
	`, bars, x, menu)
}

func BottomSpacer() string {
	return `
		<div class='mt-32'></div>
	`

}

func DynamicBanner(title string) string {
	return fmt.Sprintf(`
		<header class='flex flex-row justify-between items-center p-6 h-[75px] z-40 bg-black w-full'>
			<h2 class='text-xl'>%s</h2>
		</header>
	`, title)

}

func BreadCrumbs(breadcrumbs ...string) string {
    breadcrumbHTML := ""
	chevron := "<p class='text-gray'>»</p>"
    for index, breadcrumb := range breadcrumbs {
        breadcrumbHTML += fmt.Sprintf(`
			%s
        `, breadcrumb)
		if index < len(breadcrumbs) - 1 {
			breadcrumbHTML += chevron
		}
    }
    return fmt.Sprintf(`
        <nav class='flex flex-wrap items-center gap-4 p-6'>
            %s
        </nav>
    `, breadcrumbHTML)
}


func TitleAndText(title string, text string) string {
	return fmt.Sprintf(`
		<div class='p-6 flex gap-4 flex-col'>
			<h2 class='text-lg'>%s</h2>
			<p class='text-sm'>%s</p>
		</div>
	`, title, text)

}

func NavLink(name string, href string, isActive bool) string {
	var activeClass string
	if isActive {
		activeClass = "text-red"
	} else {
		activeClass = "text-white"
	}
	return fmt.Sprintf(`
		<a href='%s' hx-indicator='#main-loader' class='underline text-xs %s'>%s</a>
	`, href, activeClass, name)
}

func SimpleTopNav(navButtons ...string) string {
	navButtonsHTML := strings.Join(navButtons, "")
	return fmt.Sprintf(`
		<nav class='flex flex-wrap gap-4 p-6'>
			%s
		</nav>
	`, navButtonsHTML)
}

func Root(components ...string) string {
	content := strings.Join(components, "")
	return fmt.Sprintf("<main id='root' class='grid place-items-center p-4'>%s</main>", content)
}

func CenterContentWrapper(components ...string) string {
	content := strings.Join(components, "")
	return fmt.Sprintf(`
		<div class='w-full xs:w-[480px] lg:w-[640px]'>
			%s
		</div>
	`, content)

}

func NavMenuLink(name string, href string, currentPath string) string {
	activeClass := "bg-black text-white"
	if currentPath == href {
		activeClass = "bg-white text-black font-bold"
	}
	return fmt.Sprintf(`
		<li class='border border-gray hover:border-white %s rounded flex'>
			<a href='%s' class='p-4 w-full'>%s</a>
		</li>
	`, activeClass, href, name)
}

func AppNavMenu(currentPath string) string {
	return fmt.Sprintf(`
		<nav x-show='navopen' x-transition x-cloak class='fixed right-0 top-[75px] h-full w-3/5 max-w-[300px] p-6 z-40 bg-black border-l border-darkgray'>
			<ul class='flex flex-col gap-4'>
				%s%s%s
			</ul>
		</nav>
		<div x-show='navopen' @click='navopen = !navopen' class='h-full w-full fixed top-0 left-0 bg-black opacity-50 z-30'></div>
	`, NavMenuLink("Tickets", "/app/ticket", currentPath), NavMenuLink("Manufacturers", "/app/manufacturer", currentPath), NavMenuLink("Logout", "/logout", currentPath))
}

func SvgIcon(src string, size string, xattr string, xclass string) string {
	var svgSize string
	if size == "xxs" {
		svgSize = "h-2 w-2"
	}
	if size == "xs" {
		svgSize = "h-4 w-4"
	}
	if size == "sm" {
		svgSize = "h-6 w-6"
	}
	if size == "md" {
		svgSize = "h-8 w-8"
	}
	return fmt.Sprintf(`
		<div %s class='%s %s flex items-center cursor-pointer'>
			<img src='%s'/>
		</div>
	`, xattr, svgSize, xclass, src)
}

func MqWrapperCentered(content string) string {
	return fmt.Sprintf(`
		<div class='flex justify-center items-center w-full h-full'>
			%s
		</div>
	`, content)
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
			<input name='%s' type='%s' value='%s' autocomplete='off' class="flex p-1 w-full rounded-md border border-darkgray border-input bg-black px-2 py-1 text-sm ring-offset-background placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2" />
		</div>
	`, name, spanText, serverName, inputType, value)
}

func FormSubmitButton() string {
	return `
		<input type='submit' class="flex w-full row-start-5 row-end-5 col-start-1 col-end-8 bg-white text-black mt-4 items-center justify-center whitespace-nowrap rounded-md text-xs font-medium ring-offset-background transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 bg-primary text-primary-foreground hover:bg-primary/90 h-10 px-4 py-2" type="submit"/ >
	`
}

func FormDeleteButton() string {
	return `
		<input type='submit' class="flex w-full row-start-5 row-end-5 col-start-1 col-end-8 bg-red text-black mt-4 items-center justify-center whitespace-nowrap rounded-md text-xs font-medium ring-offset-background transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 bg-primary text-white text-primary-foreground hover:bg-primary/90 h-10 px-4 py-2" type="submit" value='Delete'/>
	`
}

func FormPhotoUpload() string {
	return `
		<div class='flex flex-col text-xs w-fit rounded gap-2'>
			<label>Photo</label>
			<button id='photo-upload-button' type='button' class='text-left border border-gray hover:border-lightgray p-2 rounded'>Upload Photo</button>
			<input name='photo' id='upload-input' type='file' class='hidden'/>
		</div>
		<div id='image-preview'></div>
		<script>
			document.getElementById('photo-upload-button').addEventListener('click', () => {
				document.getElementById('upload-input').click()
			})
			document.getElementById('upload-input').addEventListener('change', (e) => {
				const file = e.target.files[0];
				if (file) {
					const reader = new FileReader();
					const imgElement = document.createElement('img');
					imgElement.alt = 'Image Preview';
					reader.onload = () => {
						imgElement.src = reader.result;
						document.getElementById('image-preview').innerHTML = '';
						document.getElementById('image-preview').appendChild(imgElement);
					};
		
					reader.readAsDataURL(file);
				}
			});
		</script>
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

func FormSuccess(message string) string {
	if message == "" {
		return ""
	} else {
		return fmt.Sprintf(`
			<p class='text-green text-xs py-2'>%s</p>
		`, message)
	}

}

func FormTitle(title string) string {
	return fmt.Sprintf(`
		<h2 class='text-lg'>%s</h2>
	`, title)
}

func FormTextAreaLabel(name string, serverName string, rows int, value string) string {
	return fmt.Sprintf(`
        <div class='flex flex-col gap-2 w-full'>
            <label class="text-xs leading-none">%s</label>
            <textarea name='%s' rows='%d' autocomplete='off' spellcheck='false' class="flex p-1 w-full rounded-md border border-darkgray border-input bg-black px-2 py-1 text-sm ring-offset-background placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2">%s</textarea>
        </div>
    `, name, serverName, rows, value)
}

func FormSelectLabel(name string, serverName string, options []string, selectedValue string) string {
	var optionHTML strings.Builder
	for _, option := range options {
		selected := ""
		if option == selectedValue {
			selected = "selected"
		}
		optionHTML.WriteString(fmt.Sprintf("<option value='%s' %s>%s</option>", option, selected, option))
	}

	return fmt.Sprintf(`
        <div class='flex flex-col gap-2 w-full'>
            <label class="text-xs leading-none">%s</label>
            <select name='%s' class="flex p-1 w-full rounded-md border border-darkgray border-input bg-black px-2 py-1 text-sm ring-offset-background placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2">
                %s
            </select>
        </div>
    `, name, serverName, optionHTML.String())
}

func LinkButton(href string, text string, active bool) string {
	var activeClass string
	if active {
		activeClass = "bg-white text-black"
	} else {
		activeClass = "bg-black text-white"
	}
	return fmt.Sprintf(`
		<a href='%s' hx-indicator='#main-loader' hx-history='false' class='cursor-pointer w-fit text-center rounded py-2 px-4 text-xs border border-darkgray hover:border-white %s'>%s</a>
	`, href, activeClass, text)
}

func LoginForm(err string, username string, password string) string {
	return fmt.Sprintf(`
		<form hx-indicator='#main-loader' method='POST' class='flex flex-col p-6 gap-4 w-full'>
			%s%s%s%s%s
		</form>
	`, FormTitle("Login"), FormError(err), FormInputLabel("Username", "username", "", username), FormInputLabel("Password", "password", "password", password), FormSubmitButton())
}

func CreateManufacturerForm(err string, success string, name string, phone string, email string, submitRedirect string) string {
	return fmt.Sprintf(`
		<form hx-indicator='#main-loader' method='POST' action='/form/manufacturer/create' class='flex flex-col p-6 gap-4 w-full'>
			<input type='hidden' name='submitRedirect' value='%s'/>
			%s%s%s%s%s%s%s
		</form>
	`, submitRedirect, FormTitle("Create Manufacturers"), FormError(err), FormSuccess(success), FormInputLabel("Name", "name", "", name), FormInputLabel("Phone", "phone", "", phone), FormInputLabel("Email", "email", "", email), FormSubmitButton())
}

func ManufacturerList(manufacturers []_model.Manufacturer) string {
	manufacturerList := ""
	for _, manufacturer := range manufacturers {
		href := fmt.Sprintf("/app/manufacturer/%d", manufacturer.ID)
		manufacturerList += fmt.Sprintf(`
			<div>
				<a href='%s' hx-indicator='#main-loader' class='border text-xs border-darkgray rounded w-fit p-2 flex flex-col gap-2 hover:border-white'>%s</a>
			</div>
		`, href, manufacturer.Name)
	}
	if len(manufacturers) == 0 {
		return `
			<div id='manufacturer-list' class='p-6 w-full'>
				<p>Oh, so empty!</p>
			</div>
		`
	} else {
		return fmt.Sprintf(`
			<div id='manufacturer-list' class='p-6 flex flex-wrap gap-4'>%s</div>
		`, manufacturerList)
	}
}

func UpdateManufacturerForm(manufacturer _model.Manufacturer, err string, success string) string {
	updateAction := fmt.Sprintf("/app/manufacturer/%d/update", manufacturer.ID)
	return fmt.Sprintf(`
		<form hx-indicator='#main-loader' action='%s' method='POST' class='flex flex-col gap-4 w-full p-6'>
				%s%s%s%s%s%s%s
			<a id='manufacturer-delete-link' hx-indicator='#main-loader' href='/app/manufacturer/%d/delete' class='hidden'></a>
		</form>
	`,
		updateAction,
		FormTitle("Update Manufacturer"),
		FormError(err),
		FormSuccess(success),
		FormInputLabel("Name", "name", "", manufacturer.Name),
		FormInputLabel("Phone", "phone", "", manufacturer.Phone),
		FormInputLabel("Email", "email", "", manufacturer.Email),
		FormSubmitButton(),
	)
}

func DeleteManufacturerForm(manufacturer _model.Manufacturer, err string) string {
	deleteAction := fmt.Sprintf("/app/manufacturer/%d/delete", manufacturer.ID)
	return fmt.Sprintf(`
		<form hx-indicator='#main-loader' action='%s' method='POST' class='flex flex-col gap-4 w-full p-6'>
			%s%s%s%s
		</form>
	`, deleteAction, FormTitle("Delete Manufacturer"), FormError(err), FormInputLabel("Type '"+manufacturer.Name+"' to delete", "name", "", ""), FormDeleteButton())

}

func CreateEquipmentForm(manufacturer _model.Manufacturer, err string, success string, submitRedirect string, nickname string, serialNumber string) string {
	return fmt.Sprintf(`
		<form id='create-equipment-form' enctype='multipart/form-data' action='/form/manufacturer/%d/equipment/create' hx-indicator='#main-loader' method='POST' class='flex flex-col p-6 gap-4 w-full'>
			%s%s%s%s%s
			<div class='flex flex-col text-xs w-fit rounded gap-2'>
				<label>Photo</label>
				<button id='upload-submit' type='button' class='text-left border border-gray hover:border-lightgray p-2 rounded'>Upload Photo</button>
				<input name='photo' id='upload-input'  type='file' class='hidden'/>
			</div>
			<div id='image-preview-wrapper'></div>
			<div>%s</div>
			<input type='hidden' name='submitRedirect' value='%s'/>
		</form>
		<script>
			qs('#upload-submit').addEventListener('click', () => {
				qs('#upload-input').click();
			});
			qs('#upload-input').addEventListener('change', (uploadEvent) => {
				const file = uploadEvent.target.files[0];
				if (!file) {
					return
				}
				const reader = new FileReader();
				reader.readAsDataURL(file);
				reader.onload = (readerEvent) => {
					console.log('loaded reader event')
					const img = document.createElement('img');
					img.src = readerEvent.target.result;
					img.onload = () => {
						console.log('loaded image element')
						qs('#image-preview-wrapper').innerHTML = '';
						qs('#image-preview-wrapper').appendChild(img);
						const canvas = document.createElement('canvas');
						canvas.width = 200;
						const scaleSize = 200 / img.width;
						canvas.height = img.height * scaleSize;
						const ctx = canvas.getContext('2d');
						ctx.drawImage(img, 0, 0, canvas.width, canvas.height);
						const srcEncoded = ctx.canvas.toDataURL('image/jpeg', 0.7);
						img.src = srcEncoded;
						img.classList.add('rounded-lg');
						img.onload = null
					}
				}
			});
		</script>
	`,
		manufacturer.ID,
		FormTitle("Create Equipment"),
		FormError(err),
		FormSuccess(success),
		FormInputLabel("Nickname", "nickname", "text", nickname),
		FormInputLabel("Serial Number", "number", "text", serialNumber),
		FormSubmitButton(),
		submitRedirect,
	)
}

func HxGetLoader(href string) string {
	return fmt.Sprintf(`
		<div hx-get='%s' hx-trigger='load' hx-swap='outerHTML' class='p-6 flex items-center justify-center'>
			<div class='h-[50px] w-[50px] border-2 border-darkgray border-t-white rounded-full animate-spin'></div>
		</div>
	`, href)
}

func ManufacturerDetails(manufacturer _model.Manufacturer) string {
	return fmt.Sprintf(`
		<div class='flex flex-col p-6 gap-4'>
			<h2 class='text-lg'>Manufacturer - %s</h2>
			<div class='text-xs'>
				<p class=''>Email: %s</p>
				<p class=''>Phone: %s</p>
			</div>
		</div>
	`, manufacturer.Name, manufacturer.Email, manufacturer.Phone)
}

func EquipmentDetails(equipment _model.Equipment, manufacturer _model.Manufacturer) string {
	equipmentPhotoBase64 := base64.StdEncoding.EncodeToString(equipment.Photo)
	return fmt.Sprintf(`
		<div class='flex flex-col p-6 gap-6'>
			<h2 class='text-lg'>Equipment - %s</h2>
			<div class='flex flex-col gap-4'>
				<div class='text-xs'>
					<p class=''>Serial Number: %s</p>
					<p class=''>Manufacturer: %s</p>
				</div>
				<div id='image-preview-wrapper'>
					<img src="data:image/jpeg;base64,%s" class="w-[200px] rounded-lg" alt="%s" />
				</div>
			</div>
		</div>
	`, equipment.Nickname, equipment.SerialNumber, manufacturer.Name, equipmentPhotoBase64, equipment.Nickname)

}

func EquipmentList(equipment []_model.Equipment) string {
	equipmentList := ""
	for _, eq := range equipment {
		equipmentList += fmt.Sprintf(`
            <a href='/app/equipment/%d' hx-indicator='#main-loader' class='border place-self-center flex-grow-0 border-darkgray rounded w-fit p-2 flex flex-col gap-2 hover:border-white' style='align-self: flex-start;'>
				<img src='data:image/jpeg;base64,%s' class='w-[100px] rounded-lg' alt='%s' />
				<div>
					<p class='text-xs'>%s</p>
					<p class='text-xs'>%s</p>
				</div>
            </a>
        `, eq.ID, base64.StdEncoding.EncodeToString(eq.Photo), eq.Nickname, eq.Nickname, eq.SerialNumber)
	}
	if len(equipment) == 0 {
		return `
            <div class='p-6 w-full'>
                <p>Oh, so empty!</p>
            </div>
        `
	} else {
		return fmt.Sprintf(`
			<div class='flex flex-col p-6 gap-6'>
				<h2 class='text-lg'>Registered Equipment</h2>			
				<div class='grid grid-cols-2 md:grid-cols-3 gap-4'>
					%s
				</div>
			</div>
        `, equipmentList)
	}
}

func Footer() string {
	return `
		<div class='mt-32'></div>
		<footer class='border-t bottom-0 fixed border-gray bg-black w-full h-12 bg-black flex justify-center items-center text-xs text-gray'>
			<p>© 2024 Repairs Log</p>
		</footer>
	`
}

func ErrorMessage(err string) string {
	if err == "" {
		return ""
	} else {
		return fmt.Sprintf(`
			<p class='text-red text-xs'>%s</p>
		`, err)
	}

}

func UpdateEquipmentForm(equipment _model.Equipment, err string, success string) string {
	return fmt.Sprintf(`
		<form enctype='multipart/form-data' method='POST' hx-swap='innerHTML show:no-scoll' action='/form/equipment/%d/update' hx-indicator='#main-loader' class='gap-4 grid w-full rounded bg-black p-6'>
			<div class='flex justify-between'>
				%s
			</div>
			<div class='flex flex-col gap-4'>
				%s%s%s%s
				<div class='flex flex-col text-xs w-fit rounded gap-2'>
					<label>Photo</label>
					<button id='upload-submit' type='button' class='text-left border border-gray hover:border-lightgray p-2 rounded'>Upload Photo</button>
					<input name='photo' id='upload-input'  type='file' class='hidden'/>
				</div>
				<div>%s</div>
			</div>
		</form>
		<script>
			qs('#upload-submit').addEventListener('click', () => {
				qs('#upload-input').click();
			});
			qs('#upload-input').addEventListener('change', (uploadEvent) => {
				const file = uploadEvent.target.files[0];
				if (!file) {
					return
				}
				const reader = new FileReader();
				reader.readAsDataURL(file);
				reader.onload = (readerEvent) => {
					const img = document.createElement('img');
					img.src = readerEvent.target.result;
					img.onload = () => {
						qs('#image-preview-wrapper').innerHTML = '';
						qs('#image-preview-wrapper').appendChild(img);
						const canvas = document.createElement('canvas');
						canvas.width = 200;
						const scaleSize = 200 / img.width;
						canvas.height = img.height * scaleSize;
						const ctx = canvas.getContext('2d');
						ctx.drawImage(img, 0, 0, canvas.width, canvas.height);
						const srcEncoded = ctx.canvas.toDataURL('image/jpeg', 0.7);
						img.src = srcEncoded;
						img.classList.add('rounded-lg');
						img.onload = null
					}
				}
			});
		</script>
	`,
		equipment.ID,
		FormTitle("Update Equipment"),
		FormError(err),
		FormSuccess(success),
		FormInputLabel("Nickname", "nickname", "text", equipment.Nickname),
		FormInputLabel("Serial Number", "number", "text", equipment.SerialNumber),
		FormSubmitButton(),
	)
}

func DeleteEquipmentForm(equipment _model.Equipment, err string, success string) string {
	return fmt.Sprintf(`
		<form hx-indicator='#main-loader' action='/app/equipment/%d/delete' method='POST' class='flex p-6 flex-col gap-4 w-full'>
			%s%s%s%s%s
		</form>
	`, equipment.ID, FormTitle("Delete Equipment"), FormError(err), FormSuccess(success), FormInputLabel("Type '"+equipment.Nickname+"' to delete", "name", "", ""), FormDeleteButton())
}

func CreateTicketForm(r *http.Request, token string, action string) string {
	return fmt.Sprintf(`
		<form action='%s' method='POST' hx-indicator='#main-loader' method='POST' class='flex flex-col p-6 gap-4 w-full'>
			%s%s%s%s%s%s%s%s
			<input type='hidden' name='publicSecurityToken' value='%s'/>
		</form>
	`, action, FormTitle("Create Tickets"), FormError(r.URL.Query().Get("err")), FormSuccess(r.URL.Query().Get("success")), FormInputLabel("What is your name?", "creator", "text", r.URL.Query().Get("creator")), FormInputLabel("What item needs repaired?", "item", "", r.URL.Query().Get("item")), FormTextAreaLabel("Describe the Problem", "problem", 2, r.URL.Query().Get("problem")), FormSelectLabel("Location", "location", []string{"Southroads", "Utica"}, _util.StringWithDefault(r.URL.Query().Get("location"), "Southroads")), FormSubmitButton(), token)
}

func TicketList(tickets []_model.Ticket) string {
    ticketList := ""
    for _, ticket := range tickets {
        ticketList += fmt.Sprintf(`
            <a href='/app/ticket/%d' class='border border-darkgray rounded p-2 text-xs hover:border-white cursor-pointer'>
                <div class='flex flex-row justify-between'>
                    <h2 class='font-bold text-lg'>%s</h2>
					<div class='text-xs'>
						<p>%s</p>
                    	<p>%s</p>
					</div>
                </div>
                <p>Item: %s</p>
                <p>Problem: %s</p>
            </a>

        `, ticket.ID, ticket.Creator, ticket.Location, ticket.CreatedAt.Format("01/02/2006"), ticket.Item, ticket.Problem)
    }
    if len(tickets) == 0 {
        return `
            <div id='ticket-list' class='p-6 w-full'>
                <p>No tickets found!</p>
            </div>
        `
    } else {
        return fmt.Sprintf(`
            <div id='ticket-list' class='w-full flex flex-col p-6 gap-6'>
                %s
            </div>
        `, ticketList)
    }
}



func ActiveLink(href string, text string, isActive bool) string {
	var xclass string
	if isActive {
		xclass = "bg-white text-black"
	} else {
		xclass = "text-white bg-black"
	}
	return fmt.Sprintf(`
		<a href='%s' hx-indicator='#main-loader' hx-swap='innerHTML show:no-scroll' class='ticket-view-option p-2 border border-gray rounded hover:border-white %s'>%s</a>
	`, href, xclass, text)
}

func TicketSearchForm() string {
	return fmt.Sprintf(`
		<form action='/partial/ticketList' method='GET' class="p-6 flex flex-col gap-6">
			<input name='status' id='status-input' type='hidden' value='%s' />
			<input name='priority' id='priority-input' type='hidden' value='%s' />
			<div class='flex flex-col gap-2'>
				<p class='text-xs'>Priority Filter</p>
				<div class='flex flex-wrap gap-2'>
					<div class='priority-button px-2 py-1 bg-white text-black rounded border border-darkgray text-xs'>%s</div>
					<div class='priority-button px-2 py-1 rounded border border-darkgray text-xs'>%s</div>
					<div class='priority-button px-2 py-1 rounded border border-darkgray text-xs'>%s</div>
					<div class='priority-button px-2 py-1 rounded border border-darkgray text-xs'>%s</div>
					<div class='priority-button px-2 py-1 rounded border border-darkgray text-xs'>all</div>
				</div>
			</div>
			<div class='flex flex-col gap-2'>
			<p class='text-xs'>Status Filter</p>
			<div class='flex flex-wrap gap-2'>
				<div class='status-button px-2 py-1 bg-white text-black rounded border border-darkgray text-xs'>%s</div>
				<div class='status-button px-2 py-1 rounded border border-darkgray text-xs'>%s</div>
				<div class='status-button px-2 py-1 rounded border border-darkgray text-xs'>%s</div>
				<div class='status-button px-2 py-1 rounded border border-darkgray text-xs'>%s</div>
				<div class='status-button px-2 py-1 rounded border border-darkgray text-xs'>all</div>
			</div>
		</div>
			<div id='ticket-search' class='flex justify-between bg-black border border-darkgray p-2 rounded-lg'>
				<div class='h-6 w-6'>
					<img src="/static/svg/search-dark.svg">
				</div>
				<div class='flex flex-row justify-between items-center w-full'>
					<input name='search' id='ticket-search-input' type="text" class="pl-4 w-full bg-inherit focus:outline-none" placeholder="Search Tickets...">
					<div id='ticket-list-indicator' class='hidden animate-spin rounded-full border border-darkgray border-t-white h-6 w-6 p-2'></div>
				</div>
			</div>
		</form>
		<script>
			document.body.addEventListener = function() {};
			document.body.addEventListener('click', (event) => {
				if (qs('#ticket-search')) {
					if (!qs('#ticket-search').contains(event.target)) {
						qs('#ticket-search').classList.remove('border-lightgray');
					}
				}
			});
			qs('#ticket-search').addEventListener('click', (event) => {
				event.stopPropagation();
				qs('#ticket-search').classList.toggle('border-lightgray');
			});
			qs('#ticket-search').addEventListener('input', async (event) => {
				qs('#ticket-list-indicator').classList.remove('hidden');
				let status = qs('#status-input').value;
				let priority = qs('#priority-input').value;
				let search = event.target.value;
				let path = "/partial/ticketList?status=" + status + "&priority=" + priority + "&search=" + search;
				let res = await fetch(path, {
					method: 'GET',
				});
				if (res.status !== 200) {
				}
				let html = await res.text();
				qs('#ticket-list').outerHTML = html;
				qs('#ticket-list-indicator').classList.add('hidden');
			})
			document.querySelectorAll('.status-button').forEach(button => {
				button.addEventListener('click', async () => {
					document.querySelectorAll('.status-button').forEach(btn => {
						btn.classList.remove('bg-white', 'text-black');
						btn.classList.add('bg-black', 'text-white');
					});
					button.classList.remove('bg-black', 'text-white');
					button.classList.add('bg-white', 'text-black');
					document.getElementById('status-input').value = button.textContent.trim();
					qs('#ticket-list-indicator').classList.remove('hidden');
					let status = qs('#status-input').value;
					let priority = qs('#priority-input').value;
					let search = qs('#ticket-search-input').value;
					let path = "/partial/ticketList?status=" + status + "&priority=" + priority + "&search=" + search;
					let res = await fetch(path, {
						method: 'GET',
					});
					if (res.status !== 200) {
					}
					let html = await res.text();
					qs('#ticket-list').outerHTML = html;
					qs('#ticket-list-indicator').classList.add('hidden');
				});
			});
			document.querySelectorAll('.priority-button').forEach(button => {
				button.addEventListener('click', async () => {
					document.querySelectorAll('.priority-button').forEach(btn => {
						btn.classList.remove('bg-white', 'text-black');
						btn.classList.add('bg-black', 'text-white');
					});
					button.classList.remove('bg-black', 'text-white');
					button.classList.add('bg-white', 'text-black');
					document.getElementById('priority-input').value = button.textContent.trim();
					qs('#ticket-list-indicator').classList.remove('hidden');
					let status = qs('#status-input').value;
					let priority = qs('#priority-input').value;
					let search = qs('#ticket-search-input').value;
					let path = "/partial/ticketList?status=" + status + "&priority=" + priority + "&search=" + search;
					let res = await fetch(path, {
						method: 'GET',
					});
					if (res.status !== 200) {
					}
					let html = await res.text();
					qs('#ticket-list').outerHTML = html;
					qs('#ticket-list-indicator').classList.add('hidden');
				});
			});
		</script>
	`, 
		_model.TicketStatusNew,
		_model.TicketPriorityUnspecified,
		_model.TicketPriorityUnspecified, 
		_model.TicketPriorityLow, 
		_model.TicketPriorityInconvenient, 
		_model.TicketPriorityUrgent, 
		_model.TicketStatusNew, 
		_model.TicketStatusActive, 
		_model.TicketStatusComplete, 
		_model.TicketStatusOnHold,
	)
}

func TicketViewOptions(ticketFilter string) string {
	isNew := false
	isActive := false
	isComplete := false
	switch ticketFilter {
	case "", "new":
		isNew = true
	case "complete":
		isComplete = true
	case "active":
		isActive = true
	}

	return fmt.Sprintf(`
		<div id='ticket-view-options' class='p-6 w-full flex flex-wrap gap-4 text-xs'>
			%s%s%s
		</div>
	`, ActiveLink("/app/ticket?ticketFilter=new", "New Tickets", isNew), ActiveLink("/app/ticket?ticketFilter=active", "Active Tickets", isActive), ActiveLink("/app/ticket?ticketFilter=complete", "Completed Tickets", isComplete))
}

func TicketSettings(ticket _model.Ticket, form string) string {
	resetEquipmentButton := ""
	assignEquipmentButton := ""
	publicDetailsButton := ""
	isAssignmentForm := false
	isUpdateForm := false
	isDeleteForm := false
	isPublicForm := false
	switch form {
		case "assign":
			isAssignmentForm = true
		case "public":
			isPublicForm = true
		case "update":
			isUpdateForm = true
		case "delete":
			isDeleteForm = true
	}
	if ticket.EquipmentID == nil {
		assignEquipmentButton = LinkButton(fmt.Sprintf("/app/ticket/%d?form=assign", ticket.ID), "Assign Equipment", isAssignmentForm)
	} else {
		publicDetailsButton = LinkButton(fmt.Sprintf("/app/ticket/%d?form=public", ticket.ID), "Update Public Details", isPublicForm)
		resetEquipmentButton = LinkButton(fmt.Sprintf("/app/ticket/%d/resetEquipment", ticket.ID), "Reset Equipment Association", false)
	}
	return fmt.Sprintf(`
			<div class='p-6 w-full flex flex-col gap-4 text-xs'>
				%s%s%s%s%s
			</div>
		`,
		assignEquipmentButton,
		publicDetailsButton,
		LinkButton(fmt.Sprintf("/app/ticket/%d?form=update", ticket.ID), "Update Private Details", isUpdateForm),
		LinkButton(fmt.Sprintf("/app/ticket/%d?form=delete", ticket.ID), "Delete Ticket", isDeleteForm),
		resetEquipmentButton,
	)
		
}

func TicketDetails(ticket _model.Ticket, err string) string {
    equipmentWarning := ""
    if ticket.EquipmentID == nil {
        equipmentWarning = "<p class='text-red text-xs'>This ticket will not be made public until assigned to a piece of equipment.</p>"
    }
	missingFieldsWarning := ""
	if ticket.Owner == "" && equipmentWarning == "" || ticket.Notes == "" && equipmentWarning == "" {
		missingFieldsWarning = "<p class='text-red text-xs'>This ticket will not be made public until no fields are missing.</p>"
	}
	missingWarning := "<span class='text-red'>missing!</span>"
	var notes string
	if ticket.Notes == "" {
		notes = `<p>Notes: ` + missingWarning + `</p>`
	} else {
		notes = `<p>Notes: ` + ticket.Notes + `</p>`
	}
	var owner string
	if ticket.Owner == "" {
		owner = `<p>Owner: ` + missingWarning + `</p>`
	} else {
		owner = `<p>Owner: ` + ticket.Owner + `</p>`
	}
    return fmt.Sprintf(`
        <div class='flex flex-col p-6 gap-6'>
			<div class='flex flex-col gap-2'>
            	<h2 class='text-lg'>Ticket Details</h2>
        	    %s%s
			</div>
            <div class='flex flex-col gap-1 text-xs'>
                <p>Creator: %s</p>
                <p>Created: %s</p>
                <p>Item: %s</p>
                <p>Problem: %s</p>
                <p>Location: %s</p>
                <p>Priority: %s</p>
                <p>Status: %s</p>
                %s
                %s
            </div>
        </div>
    `, equipmentWarning, missingFieldsWarning, ticket.Creator, ticket.CreatedAt.Format("01/02/2006"), ticket.Item, ticket.Problem, ticket.Location, string(ticket.Priority), string(ticket.Status), owner, notes)
}

func DeleteTicketForm(ticket _model.Ticket, err string) string {
	return fmt.Sprintf(`
		<form id='delete-ticket-form' hx-indicator='#main-loader' action='/app/ticket/%d/delete' method='POST' class='flex flex-col gap-4 w-full p-6'>
			%s%s%s%s
		</form>
	`, ticket.ID, FormTitle("Delete Ticket"), FormError(err), FormInputLabel("Type 'yes' to delete", "keyword", "", ""), FormDeleteButton())
}

func TicketAssignmentForm(ticket _model.Ticket, manufacturers []_model.Manufacturer, db *gorm.DB) string {
	manufacturerOptions := ""
	for _, manufacturer := range manufacturers {
		manufacturerOptions += fmt.Sprintf("<div hx-get='/partial/manufacturer/%d/equipmentSelectionList?ticketID=%d' hx-indicator='#main-loader' hx-target='#equipment-selection-list' hx-swap='outerHTML' class='manufacturer-option border h-fit hover:border-white flex items-center justify-center border-darkgray cursor-pointer bg-black p-2 rounded text-xs' value='%d'>%s</div>", manufacturer.ID, ticket.ID, manufacturer.ID, manufacturer.Name)
	}
	return fmt.Sprintf(`
		<form id='ticket-assignment-form' action='/app/ticket/%d/assign' hx-indicator='#main-loader' method='POST' class='flex flex-col gap-6 w-full p-6'>
			<input id='equipment-id-input' type='hidden' name='equipmentID' value=''/>
			<div id='manufacturer-selection-list' class='flex flex-col gap-6'>
				<div class='flex flex-row justify-between items-center'>
					<h2>Assign Manufacturer</h2>
				</div>
				<div class='flex flex-wrap gap-4 text-xs'>
					%s
				</div>
				<div class='hidden'>
					<input id='equipment-selection-hidden-submit' type='submit' class='hidden'/>
				</div>
			</div>
			</form>
			<div id='equipment-selection-list' class='hidden'></div>
		<span id='ticket-assignment-submit' class='hidden flex p-6'>
			<input type='submit' class='flex w-full bg-white text-black mt-4 items-center justify-center whitespace-nowrap rounded-md text-xs font-medium ring-offset-background transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 bg-primary text-primary-foreground hover:bg-primary/90 h-10 px-4 py-2' type='submit' value='Assign'/>
		</span>
		<script>
		document.querySelectorAll('.manufacturer-option').forEach((option) => {
			option.addEventListener('click', (e) => {
				document.querySelectorAll('.manufacturer-option').forEach((opt) => {
					opt.classList.remove('bg-white', 'text-black');
					opt.classList.add('bg-black', 'text-white');
				});
				e.target.classList.remove('bg-black', 'text-white');
				e.target.classList.add('bg-white', 'text-black');
				if (document.getElementById('ticket-assignment-submit').classList.contains('hidden')) {
				
				} else {
					document.getElementById('ticket-assignment-submit').classList.add('hidden');
				}
			});
		});
		qs("#ticket-assignment-submit").addEventListener('click', () => {
			document.getElementById('equipment-selection-hidden-submit').click();
		})
		</script>
	`, ticket.ID, manufacturerOptions)
}

func MainLoader() string {
    return `
        <div id='main-loader' class='flex-indicator'>
			<div class='fixed top-0 left-0 w-full h-full bg-black opacity-50 z-50'></div>
            <div class='fixed top-1/3 left-1/2 transform -translate-x-1/2 -translate-y-1/2 z-50'>
                <div class='h-[100px] w-[100px] rounded-full border-4 border-darkgray border-t-white animate-spin'></div>
            </div>
        </div>
    `
}

func ResetEquipmentLink(ticket _model.Ticket) string {
	return fmt.Sprintf(`
		<div class='p-6 flex flex-col'>
			<h2>Equipment Association</h2>
			<p>This ticket is associated with the following piece of equipment:</p>
			
		</div>
	`, ticket.ID)
}

func UpdateTicketForm(ticket _model.Ticket, err string, success string) string {
	return fmt.Sprintf(`
		<form id='combined-form' hx-indicator='#main-loader' hx-swap='innerHTML scroll:no-show' action='/form/ticket/%d/update' method='POST' class='flex flex-col gap-4 w-full p-6'>
			<div class='flex justify-between'>
				%s
			</div>
			<div class='flex flex-col gap-4'>
				%s%s%s%s%s%s%s%s%s%s%s
			</div>
		</form>
	`,
		ticket.ID,
		FormTitle("Update Ticket Details"),
		FormError(err),
		FormSuccess(success),
		FormInputLabel("Creator", "creator", "text", ticket.Creator),
		FormInputLabel("Item Description", "item", "text", ticket.Item),
		FormTextAreaLabel("Problem", "problem", 2, ticket.Problem),
		FormSelectLabel("Location", "location", []string{"Southroads", "Utica"}, ticket.Location),
		FormInputLabel("Owner", "owner", "text", ticket.Owner),
		FormSelectLabel("Priority", "priority", []string{string(_model.TicketPriorityUnspecified), string(_model.TicketPriorityLow), string(_model.TicketPriorityInconvenient), string(_model.TicketPriorityUrgent)}, string(ticket.Priority)),
		FormSelectLabel("Status", "status", []string{string(_model.TicketStatusNew), string(_model.TicketStatusActive), string(_model.TicketStatusOnHold), string(_model.TicketStatusComplete)}, string(ticket.Status)),
		FormTextAreaLabel("Notes", "notes", 2, ticket.Notes),
		FormSubmitButton(),
	)
}