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
		<a href='%s' hx-indicator='#main-loader' class='cursor-pointer w-fit text-center rounded py-1 px-4 border border-white %s'>%s</a>
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
		<form hx-indicator='#main-loader' method='POST' action='/app/manufacturer' class='flex flex-col p-6 gap-4 w-full'>
			<input type='hidden' name='submitRedirect' value='%s'/>
			%s%s%s%s%s%s%s
		</form>
	`, submitRedirect, FormTitle("Create Manufacturer"), FormError(err), FormSuccess(success), FormInputLabel("Name", "name", "", name), FormInputLabel("Phone", "phone", "", phone), FormInputLabel("Email", "email", "", email), FormSubmitButton())
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
			<div class='p-6 w-full %s'>
				<p>Oh, so empty!</p>
			</div>
		`, xclass)
	} else {
		return fmt.Sprintf(`
			<div class='p-6 w-full %s'>
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

func UpdateManufacturerForm(manufacturer _model.Manufacturer, err string, success string) string {
	updateAction := fmt.Sprintf("/app/manufacturer/%d/update", manufacturer.ID)
	return fmt.Sprintf(`
		<form hx-indicator='#main-loader' action='%s' method='POST' class='flex flex-col gap-4 w-full p-6'>
			%s%s%s%s%s%s%s
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

func ManufacturerDetails(manufacturer _model.Manufacturer, err string, form string) string {
	isCreateEquipment := false
	isUpdateManufacturer := false
	isDeleteManufacturer := false
	switch form {
		case "create", "":
			isCreateEquipment = true

		case "update":
			isUpdateManufacturer = true
		case "delete":
			isDeleteManufacturer = true
	}

	return fmt.Sprintf(`
		<div class='p-6 w-full flex flex-col gap-6'>
			<div class='flex flex-row justify-between'>
				<h2>Manufacturer Details</h2>
			</div>
			<div class='text-xs'>
				<p><strong>Name:</strong> %s</p>
				<p><strong>Phone:</strong> %s</p>
				<p><strong>Email:</strong> %s</p>
			</div>
			<div id='hidden-settings-section' class='text-xs flex flex-col gap-4'>
				%s%s%s
			</div>
		</div>
	`,
		manufacturer.Name,
		manufacturer.Phone,
		manufacturer.Email,
		LinkButton(fmt.Sprintf("/app/manufacturer/%d?form=create", manufacturer.ID), "Create Equipment", isCreateEquipment),
		LinkButton(fmt.Sprintf("/app/manufacturer/%d?form=update", manufacturer.ID), "Update Manufacturer", isUpdateManufacturer),
		LinkButton(fmt.Sprintf("/app/manufacturer/%d?form=delete", manufacturer.ID), "Delete Manufacturer", isDeleteManufacturer),
	)
}

func CreateEquipmentForm(manufacturer _model.Manufacturer, err string, success string, submitRedirect string, nickname string, serialNumber string) string {
	return fmt.Sprintf(`
		<form enctype='multipart/form-data' action='/app/manufacturer/%d' hx-indicator='#main-loader' method='POST' class='flex flex-col p-6 gap-4 w-full'>
			%s%s%s%s%s
			<div class='flex flex-col text-xs w-fit rounded gap-2'>
				<label>Photo</label>
				<button id='upload-submit' type='button' class='text-left border border-gray hover:border-lightgray p-2 rounded'>Upload Photo</button>
				<input name='photo' id='upload-input'  type='file' class='hidden'/>
			</div>
			<div id='image-preview'></div>
			%s
			<input type='hidden' name='submitRedirect' value='%s'/>
		</form>
		<script>
			qs('#upload-submit').addEventListener('click', () => {
				qs('#upload-input').click()
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

func EquipmentList(equipment []_model.Equipment) string {
	equipmentList := ""
	for _, eq := range equipment {
		equipmentList += fmt.Sprintf(`
            <a href='/app/equipment/%d' hx-indicator='#main-loader' class='border flex-grow-0 border-darkgray rounded w-fit p-2 flex flex-col gap-2 hover:border-white' style='align-self: flex-start;'>
				<img src='data:image/jpeg;base64,%s' class='w-[100px] rounded-lg' alt='%s' />
				<p class='text-sm'>%s</p>
            </a>
        `, eq.ID, base64.StdEncoding.EncodeToString(eq.Photo), eq.Nickname, eq.Nickname)
	}
	if len(equipment) == 0 {
		return `
            <div class='p-6 w-full'>
                <p>Oh, so empty!</p>
            </div>
        `
	} else {
		return fmt.Sprintf(`
            <div class='p-6 w-full flex flex-wrap gap-6'>
				%s
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

func EquipmentDetails(equipment _model.Equipment, manufacturer _model.Manufacturer, form string) string {
	isUpdateForm := false
	isDeleteForm := false
	switch form {
		case "update":
			isUpdateForm = true
		case "delete":
			isDeleteForm = true
	}
	return fmt.Sprintf(`
		<div class='p-6 w-full grid gap-6'>
			<h2>Equipment Details</h2>
			<div class='text-xs'>
				<p class=''><strong>Nickname:</strong> %s</p>
				<p class=''><strong>Serial Number:</strong> %s</p>
				<p class=''><strong>Manufacturer:</strong> <a href='/app/manufacturer/%d' class='underline hover:text-red'>%s</a></p>
			</div>
			<div class='flex flex-row justify-between'>
				<img src='data:image/jpeg;base64,%s' class='w-[200px] rounded-lg' alt='%s'/>
			</div>
			<div id='hidden-settings-section' class='text-xs flex flex-col gap-4 mt-2'>
				%s%s
			</div>
		</div>
	`,
		equipment.Nickname,
		equipment.SerialNumber,
		manufacturer.ID,
		manufacturer.Name,
		base64.StdEncoding.EncodeToString(equipment.Photo),
		equipment.Nickname,
		LinkButton(fmt.Sprintf("/app/equipment/%d?form=update", equipment.ID), "Update Equipment", isUpdateForm),
		LinkButton(fmt.Sprintf("/app/equipment/%d?form=delete", equipment.ID), "Delete Equipment", isDeleteForm),
	)
}

func UpdateEquipmentForm(equipment _model.Equipment, err string, success string) string {
	return fmt.Sprintf(`
		<form enctype='multipart/form-data' method='POST' hx-swap='innerHTML show:no-scoll' action='/app/equipment/%d/update' hx-indicator='#main-loader' class='gap-4 grid w-full rounded bg-black p-6'>
			<div class='flex justify-between'>
				%s
			</div>
			<div class='flex flex-col gap-4'>
				%s%s%s%s%s%s
			</div>
		</form>
	`,
		equipment.ID,
		FormTitle("Update Equipment"),
		FormError(err),
		FormSuccess(success),
		FormInputLabel("Nickname", "nickname", "text", equipment.Nickname),
		FormInputLabel("Serial Number", "number", "text", equipment.SerialNumber),
		FormPhotoUpload(),
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

func EquipmentQrCodeDownload(equipment _model.Equipment) string {
	return fmt.Sprintf(`
		<div class='flex flex-row gap-4 justify-between p-6'>
			<h2>Equipment QR Code</h2>
			<div id='qrcode-download-icon' href='/app/equipment/%d/downloadticketqr'>%s</div>
		</div>
		<script>
			document.getElementById('qrcode-download-icon').addEventListener('click', () => {
				window.location.href = '/app/equipment/%d/downloadticketqr'
			})
		</script>
	`, equipment.ID, SvgIcon("/static/svg/download-dark.svg", "sm", "", ""), equipment.ID)
}

func CreateTicketForm(r *http.Request, token string) string {
	return fmt.Sprintf(`
		<form enctype='multipart/form-data' action='/app/ticket/public' method='POST' hx-indicator='#main-loader' method='POST' class='flex flex-col p-6 gap-4 w-full'>
			%s%s%s%s%s%s%s%s%s
			<input type='hidden' name='publicSecurityToken' value='%s'/>
		</form>
	`, FormTitle("Create Tickets"), FormError(r.URL.Query().Get("err")), FormSuccess(r.URL.Query().Get("success")), FormInputLabel("What is your name?", "creator", "text", r.URL.Query().Get("creator")), FormInputLabel("What item needs repaired?", "item", "", r.URL.Query().Get("item")), FormTextAreaLabel("Describe the Problem", "problem", 2, r.URL.Query().Get("problem")), FormSelectLabel("Location", "location", []string{"Southroads", "Utica"}, _util.StringWithDefault(r.URL.Query().Get("location"), "Southroads")), FormPhotoUpload(), FormSubmitButton(), token)
}

func TicketList(tickets []_model.Ticket) string {
	ticketList := ""
	for _, ticket := range tickets {
		ticketList += fmt.Sprintf(`
			<a href='/app/ticket/%d' class='flex flex-col border rounded border-gray text-xs'>
				<div class='flex border-b border-gray p-2 rounded-t'>
					<div class='font-bold text-sm'>%s</div>			
				</div>
				<div class='p-2'>
					<div class=''>Creator: %s</div>
					<div class=''>Location: %s</div>
				</div>
				<div class='p-2'>Problem: %s</div>
			</a>

        `, ticket.ID, ticket.Item, ticket.Creator, ticket.Location, ticket.Problem)
	}
	if len(tickets) == 0 {
		return `
            <div class='p-6 w-full'>
                <p>No tickets found!</p>
            </div>
        `
	} else {
		return fmt.Sprintf(`
            <div class='p-6 w-full flex flex-col text-sm gap-6'>
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
		<a href='%s' hx-swap='innerHTML show:no-scroll' class='ticket-view-option p-2 border border-gray rounded hover:border-white %s'>%s</a>
	`, href, xclass, text)
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

func TicketSettings(ticket _model.Ticket) string {
	resetEquipmentButton := ""
	if ticket.EquipmentID != nil {
		resetEquipmentButton = fmt.Sprintf(`<a id='ticket-reset-equipment-button' href='/app/ticket/%d/resetEquipment' hx-indicator='#main-loader' class='cursor-pointer text-center rounded w-fit py-1 px-4 border border-black bg-green text-black text-sm'>Reset Equipment Association</a>`, ticket.ID)
	}
	assignEquipmentButton := ""
	if ticket.EquipmentID == nil {
		assignEquipmentButton = fmt.Sprintf(`<a href='/app/ticket/%d?form=assign' hx-indicator='#main-loader' class='cursor-pointer text-center rounded w-fit py-1 px-4 border border-black bg-blue text-white text-sm'>Assign Equipment</a>`, ticket.ID)
	}
	return fmt.Sprintf(`
			<div id='ticket-settings' class='text-sm flex flex-col gap-4 hidden p-6'>
				<a href='/app/ticket/%d?form=update' hx-indicator='#main-loader' id='ticket-update-button' class='cursor-pointer w-fit text-center rounded py-1 px-4 border border-black text-black bg-white'>Update Private Details</a>
				%s
				%s
				<a href='/app/ticket/%d?form=delete' hx-indicator='#main-loader' id='ticket-delete-button' class='cursor-pointer w-fit text-center rounded py-1 px-4 border border-black bg-red'>Delete Ticket</a>				
			</div>
		`,
		ticket.ID, resetEquipmentButton, assignEquipmentButton, ticket.ID,
	)
}

func TicketDetails(ticket _model.Ticket, err string) string {
	ticketPhotoBase64 := base64.StdEncoding.EncodeToString(ticket.Photo)
	return fmt.Sprintf(`
			<div class='grid grid-cols-2 p-6 gap-6'>
				<img src='data:image/jpeg;base64,%s' class='w-full h-auto' alt='%s'/>
				<div class='flex flex-col gap-6'>
					<div id='ticket-settings-icon' class='place-self-end'>%s</div>
					<div id='ticket-close-settings-icon' class='place-self-end hidden'>%s</div>
				</div>
			</div>
		`,
		ticketPhotoBase64,
		ticket.Item,
		SvgIcon("/static/svg/gear-dark.svg", "sm", "", ""),
		SvgIcon("/static/svg/x-dark.svg", "sm", "", ""),
	)
}

func UpdateTicketForm(ticket _model.Ticket) string {
	return fmt.Sprintf(`
		<form id='update-ticket-form' hx-indicator='#main-loader' enctype='multipart/form-data' method='POST' action='/app/ticket/%d/update' class='gap-4 grid w-full rounded bg-black p-6'>
			<div class='flex justify-between'>
				%s
			</div>
			<div class='flex flex-col gap-4'>
				%s%s%s%s%s
			</div>
		</form>
	`,
		ticket.ID,
		FormTitle("Update Ticket"),
		FormInputLabel("Creator", "creator", "text", ticket.Creator),
		FormInputLabel("Item Description", "item", "text", ticket.Item),
		FormTextAreaLabel("Problem", "problem", 2, ticket.Problem),
		FormPhotoUpload(),
		FormSubmitButton(),
	)
}

func DeleteTicketForm(ticket _model.Ticket) string {
	return fmt.Sprintf(`
		<form id='delete-ticket-form' hx-indicator='#main-loader' action='/app/ticket/%d/delete' method='POST' class='flex flex-col gap-4 w-full p-6'>
			%s%s%s
		</form>
	`, ticket.ID, FormTitle("Delete Ticket"), FormInputLabel("Type 'yes' to delete", "keyword", "", ""), FormDeleteButton())
}

func TicketPublicDetailsForm(ticket _model.Ticket) string {
	return fmt.Sprintf(`
		<form id='ticket-public-details' action='/app/ticket/%d/publicdetails' method='POST' class='flex flex-col gap-4 w-full p-6'>
			%s%s%s%s%s%s
		</form>
	`,
		ticket.ID,
		FormTitle("Ticket Public Details"),
		FormInputLabel("Owner", "owner", "text", ticket.Owner),
		FormSelectLabel("Priority", "priority", []string{string(_model.TicketPriorityUnspecified), string(_model.TicketPriorityLow), string(_model.TicketPriorityInconvenient), string(_model.TicketPriorityUrgent)}, string(ticket.Priority)),
		FormSelectLabel("Status", "status", []string{string(_model.TicketStatusNew), string(_model.TicketStatusActive), string(_model.TicketStatusOnHold), string(_model.TicketStatusComplete)}, string(ticket.Status)),
		FormTextAreaLabel("Notes", "notes", 2, ticket.Notes),
		FormSubmitButton(),
	)
}

func TicketSettingsScript() string {
	return `
		<script>
			qs('#ticket-settings-icon').addEventListener('click', (e) => {
				qs('#ticket-settings').classList.toggle('hidden')
				qs('#ticket-settings-icon').classList.toggle('hidden')
				qs('#ticket-close-settings-icon').classList.toggle('hidden')
			})
			qs('#ticket-close-settings-icon').addEventListener('click', (e) => {
				qs('#ticket-settings').classList.toggle('hidden')
				qs('#ticket-settings-icon').classList.toggle('hidden')
				qs('#ticket-close-settings-icon').classList.toggle('hidden')
			})
		</script>
	`
}

func TicketAssignmentForm(ticket _model.Ticket, manufacturers []_model.Manufacturer, db *gorm.DB) string {
	manufacturerOptions := fmt.Sprintf(`<a href='/app/manufacturer?submitRedirect=/app/ticket/%d' class='border hover:border-white border-darkgray cursor-pointer bg-black p-2 rounded-full text-sm'><img class='h-[25px] w-[25px]' src='/static/svg/plus-dark.svg'></img></a>`, ticket.ID)
	for _, manufacturer := range manufacturers {
		manufacturerOptions += fmt.Sprintf("<div hx-get='/partial/manufacturer/%d/equipmentSelectionList?ticketID=%d' hx-indicator='#main-loader' hx-target='#equipment-selection-list' hx-swap='outerHTML' class='manufacturer-option border h-fit hover:border-white flex items-center justify-center border-darkgray cursor-pointer bg-black p-2 rounded text-sm' value='%d'>%s</div>", manufacturer.ID, ticket.ID, manufacturer.ID, manufacturer.Name)
	}
	return fmt.Sprintf(`
		<form id='ticket-assignment-form' action='/app/ticket/%d/assign' hx-indicator='#main-loader' method='POST' class='flex flex-col gap-4 w-full p-6'>
			<input id='equipment-id-input' type='hidden' name='equipmentID' value=''/>
			<div id='manufacturer-selection-list' class='flex flex-col gap-6'>
				<h2>Assign Manufacturer</h2>
				<div class='flex flex-wrap gap-4'>
					%s
				</div>
			</div>
			<div id='equipment-selection-list' class='hidden'></div>
			<span id='ticket-assignment-submit' class='hidden'>%s</span>
		</form>
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
		</script>
	`, ticket.ID, manufacturerOptions, FormSubmitButton())
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
