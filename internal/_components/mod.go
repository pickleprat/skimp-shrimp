package _components

import (
	"cfasuite/internal/_model"
	"cfasuite/internal/_util"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
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

func FormLoader() string {
	return `
		<div x-show='loading' x-cloak class='flex items-center justify-center mt-4'>
			<div class='h-10 w-10 rounded-full border-2 border-gray border-t-white animate-spin'></div>
		</div>
	`
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

func LoginForm(err string, username string, password string) string {
	return fmt.Sprintf(`
		<form x-data="{ loading: false }" method='POST' class='flex flex-col p-6 gap-4 w-full'>
			%s%s%s%s%s%s
		</form>
	`, FormTitle("Login"), FormError(err), FormInputLabel("Username", "username", "", username), FormInputLabel("Password", "password", "password", password), FormLoader(), FormSubmitButton())
}

func CreateManufacturerForm(err string, name string, phone string, email string, xclass string) string {
	return fmt.Sprintf(`
		<form x-data="{ loading: false }" method='POST' action='/app/manufacturer' class='flex flex-col p-6 gap-4 w-full %s'>
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

func UpdateManufacturerForm(manufacturer _model.Manufacturer, err string) string {
	updateAction := fmt.Sprintf("/app/manufacturer/%d/update", manufacturer.ID)
	return fmt.Sprintf(`
		<form x-data="{ loading: false }" action='%s' method='POST' class='flex flex-col gap-4 w-full'>
			%s%s%s%s%s%s%s
		</form>
	`,
		updateAction,
		FormTitle("Update Manufacturer"),
		FormError(err),
		FormInputLabel("Name", "name", "", manufacturer.Name),
		FormInputLabel("Phone", "phone", "", manufacturer.Phone),
		FormInputLabel("Email", "email", "", manufacturer.Email),
		FormLoader(),
		FormSubmitButton(),
	)
}

func DeleteManufacturerForm(manufacturer _model.Manufacturer, err string) string {
	deleteAction := fmt.Sprintf("/app/manufacturer/%d/delete", manufacturer.ID)
	return fmt.Sprintf(`
		<form x-data="{ loading: false }" action='%s' method='POST' class='flex flex-col gap-4 w-full'>
			%s%s%s%s
		</form>
	`, deleteAction, FormTitle("Delete Manufacturer"), FormError(err), FormInputLabel("Type '"+manufacturer.Name+"' to delete", "name", "", ""), FormDeleteButton())

}

func ManufacturerDetails(manufacturer _model.Manufacturer, err string) string {
	return fmt.Sprintf(`
		<div class='p-6 w-full flex flex-col gap-4'>
			<div class='flex flex-row justify-between'>
				<h2>Manufacturer Details</h2>
				%s%s
			</div>
			<div class='text-xs'>
				<p><strong>Name:</strong> %s</p>
				<p><strong>Phone:</strong> %s</p>
				<p><strong>Email:</strong> %s</p>
			</div>
			<div id='hidden-settings-section' class='text-xs hidden flex flex-row gap-2'>
				<div id='manufacturer-update-button' class='cursor-pointer rounded py-1 px-2 border border-black bg-white text-black'>Update</div>
				<div id='manufacturer-delete-button' class='cursor-pointer rounded py-1 px-2 border border-black bg-red'>Delete</div>
			</div>
			<div id='manufacturer-delete-message'>
				%s
			</div>
			<div id='update-manufacturer-form' class='hidden mt-2'>
				%s
			</div>
			<div id='delete-manufacturer-form' class='hidden mt-2'>
				%s
			</div>
		</div>
		<script>
			document.getElementById('manufacturer-settings-icon').addEventListener('click', () => {
				document.getElementById('hidden-settings-section').classList.toggle('hidden')
				document.getElementById('manufacturer-close-settings-icon').classList.toggle('hidden')
				document.getElementById('manufacturer-settings-icon').classList.toggle('hidden')
				document.getElementById('manufacturer-delete-message').classList.add('hidden')
			})
			document.getElementById('manufacturer-close-settings-icon').addEventListener('click', () => {
				document.getElementById('hidden-settings-section').classList.toggle('hidden')
				document.getElementById('manufacturer-close-settings-icon').classList.toggle('hidden')
				document.getElementById('manufacturer-settings-icon').classList.toggle('hidden')
				document.getElementById('delete-manufacturer-form').classList.add('hidden')
				document.getElementById('update-manufacturer-form').classList.add('hidden')
			})
			document.getElementById('manufacturer-update-button').addEventListener('click', () => {
				document.getElementById('update-manufacturer-form').classList.toggle('hidden')
				document.getElementById('delete-manufacturer-form').classList.add('hidden')
			})
			document.getElementById('manufacturer-delete-button').addEventListener('click', () => {
				document.getElementById('delete-manufacturer-form').classList.toggle('hidden')
				document.getElementById('update-manufacturer-form').classList.add('hidden')
			})

		</script>
	`,
		SvgIcon("/static/svg/gear-dark.svg", "sm", "id='manufacturer-settings-icon'", ""),
		SvgIcon("/static/svg/x-dark.svg", "sm", "id='manufacturer-close-settings-icon'", "hidden"),
		manufacturer.Name,
		manufacturer.Phone,
		manufacturer.Email,
		ErrorMessage(err),
		UpdateManufacturerForm(manufacturer, ""),
		DeleteManufacturerForm(manufacturer, ""),
	)
}

func CreateEquipmentForm(err string, xclass string) string {
	return fmt.Sprintf(`
		<form enctype='multipart/form-data' x-data="{ loading: false }" method='POST' class='flex flex-col p-6 gap-4 w-full  %s'>
			%s%s%s%s
			<div class='flex flex-col text-xs w-fit rounded gap-2'>
				<label>Photo</label>
				<button type='button' x-on:click="$refs.uploadInput.click()" class='text-left border border-gray hover:border-lightgray p-2 rounded'>Upload Photo</button>
				<input name='photo' id='upload-input' x-ref='uploadInput' type='file' class='hidden'/>
			</div>
			<div id='image-preview'></div>
			%s%s
		</form>
		<script>
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
		xclass,
		FormTitle("Create Equipment"),
		FormError(err),
		FormInputLabel("Nickname", "nickname", "", ""),
		FormInputLabel("Serial Number", "number", "", ""),
		FormLoader(),
		FormSubmitButton(),
	)
}

func EquipmentList(equipment []_model.Equipment, xclass string) string {
	equipmentList := ""
	for _, eq := range equipment {
		equipmentList += fmt.Sprintf(`
            <tr>
                <td>
                    <a href='/app/equipment/%d' class='underline select-none hover:text-red cursor-pointer'>%s</a>
                </td>
                <td>%s</td>
            </tr>
        `, eq.ID, eq.Nickname, eq.SerialNumber)
	}
	if len(equipment) == 0 {
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
                        <th>Nickname</th>
                        <th>Serial Number</th>
                    </tr>
                    %s
                </table>
            </div>
        `, xclass, equipmentList)
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

func EquipmentDetails(equipment _model.Equipment, manufacturer _model.Manufacturer, err string) string {
	return fmt.Sprintf(`
		<div class='p-6 w-full grid gap-4'>
			<div class='flex flex-row justify-between'>
				<h2 class='mb-2 text-lg'>Equipment Details</h2>
				%s%s
			</div>
			<div id='equipment-crud-err' class=''>
				%s
			</div>
			<div>
				<p class='text-xs'>Nickname: %s</p>
				<p class='text-xs'>Serial Number: %s</p>
				<p class='text-xs'>Manufacturer: <a href='/app/manufacturer/%d' class='underline hover:text-red'>%s</a></p>
			</div>
			<div class='w-[200px]'>
				<img src='data:image/jpeg;base64,%s' class='w-full h-auto' alt='%s'/>
			</div>
			<div id='hidden-settings-section' class='text-xs hidden flex flex-row gap-2 mt-2'>
				<div id='manufacturer-update-button' class='cursor-pointer rounded py-1 px-2 border border-black bg-white text-black'>Update</div>
				<div id='manufacturer-delete-button' class='cursor-pointer rounded py-1 px-2 border border-black bg-red'>Delete</div>
			</div>
			<div id='update-equipment-form' class='mt-2 hidden'>
				%s
			</div>
			<div id='delete-equipment-form' class='mt-2 hidden'>
				%s
			</div>
		</div>
		<script>
			document.getElementById('equipment-settings-icon').addEventListener('click', () => {
				document.getElementById('hidden-settings-section').classList.toggle('hidden')
				document.getElementById('equipment-settings-icon').classList.toggle('hidden')
				document.getElementById('equipment-close-settings-icon').classList.toggle('hidden')
				document.getElementById('equipment-crud-err').classList.add('hidden')
			})
			document.getElementById('equipment-close-settings-icon').addEventListener('click', () => {
				document.getElementById('hidden-settings-section').classList.toggle('hidden')
				document.getElementById('equipment-settings-icon').classList.toggle('hidden')
				document.getElementById('equipment-close-settings-icon').classList.toggle('hidden')
				document.getElementById('update-equipment-form').classList.add('hidden')
				document.getElementById('delete-equipment-form').classList.add('hidden')
			})
			document.getElementById('manufacturer-update-button').addEventListener('click', () => {
				document.getElementById('update-equipment-form').classList.remove('hidden')
				document.getElementById('delete-equipment-form').classList.add('hidden')
			})
			document.getElementById('manufacturer-delete-button').addEventListener('click', () => {
				document.getElementById('delete-equipment-form').classList.remove('hidden')
				document.getElementById('update-equipment-form').classList.add('hidden')
			})
		</script>
	`,
		SvgIcon("/static/svg/gear-dark.svg", "sm", "id='equipment-settings-icon'", ""),
		SvgIcon("/static/svg/x-dark.svg", "sm", "id='equipment-close-settings-icon'", "hidden"),
		ErrorMessage(err),
		equipment.Nickname,
		equipment.SerialNumber,
		manufacturer.ID,
		manufacturer.Name,
		base64.StdEncoding.EncodeToString(equipment.Photo),
		equipment.Nickname,
		UpdateEquipmentForm(equipment),
		DeleteEquipmentForm(equipment),
	)
}

func UpdateEquipmentForm(equipment _model.Equipment) string {
	return fmt.Sprintf(`
		<form enctype='multipart/form-data' method='POST' action='/app/equipment/%d/update' x-data="{ loading: false }" class='gap-4 grid w-full rounded bg-black'>
			<div class='flex justify-between'>
				%s
			</div>
			<div class='flex flex-col gap-4'>
				%s%s%s%s%s
			</div>
		</form>
	`,
		equipment.ID,
		FormTitle("Update Equipment"),
		FormInputLabel("Nickname", "nickname", "text", equipment.Nickname),
		FormInputLabel("Serial Number", "number", "text", equipment.SerialNumber),
		FormPhotoUpload(),
		FormLoader(),
		FormSubmitButton(),
	)
}

func DeleteEquipmentForm(equipment _model.Equipment) string {
	return fmt.Sprintf(`
		<form x-data="{ loading: false }" action='/app/equipment/%d/delete' method='POST' class='flex flex-col gap-4 w-full'>
			%s%s%s
		</form>
	`, equipment.ID, FormTitle("Delete Equipment"), FormInputLabel("Type '"+equipment.Nickname+"' to delete", "name", "", ""), FormDeleteButton())
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
		<form enctype='multipart/form-data' action='/app/ticket/public' method='POST' x-data="{ loading: false }" method='POST' class='flex flex-col p-6 gap-4 w-full'>
			%s%s%s%s%s%s%s%s%s%s
			<input type='hidden' name='publicSecurityToken' value='%s'/>
		</form>
	`, FormTitle("Create Tickets"), FormError(r.URL.Query().Get("err")), FormSuccess(r.URL.Query().Get("success")), FormInputLabel("What is your name?", "creator", "text", r.URL.Query().Get("creator")), FormInputLabel("What item needs repaired?", "item", "", r.URL.Query().Get("item")), FormTextAreaLabel("Describe the Problem", "problem", 2, r.URL.Query().Get("problem")), FormSelectLabel("Location", "location", []string{"Southroads", "Utica"}, _util.StringWithDefault(r.URL.Query().Get("location"), "Southroads")), FormPhotoUpload(), FormSubmitButton(), FormLoader(), token)
}

func TicketList(tickets []_model.Ticket) string {
	ticketList := ""
	for _, ticket := range tickets {
		ticketList += fmt.Sprintf(`
			<a href='/app/ticket/%d' class='grid grid-cols-4 text-xs'>
				<div class='border border-gray p-1'>%s</div>
				<div class='border border-gray p-1'>%s</div>
				<div class='border border-gray p-1'>%s</div>
				<div class='border border-gray p-1'>%s</div>
			</a>

        `, ticket.ID, ticket.Creator, ticket.Item, ticket.Location, ticket.Problem)
	}
	if len(tickets) == 0 {
		return `
            <div class='p-6 w-full'>
                <p>No tickets found!</p>
            </div>
        `
	} else {
		return fmt.Sprintf(`
            <div class='p-6 w-full flex flex-col text-sm'>
                <div class='grid grid-cols-4'>
                    <div class='border p-1'>Creator</div>
                    <div class='border p-1'>Item</div>
                    <div class='border p-1'>Location</div>
                    <div class='border p-1'>Problem</div>
                </div>
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

func TicketDetails(ticket _model.Ticket, err string) string {
	return fmt.Sprintf(`
		<div class='p-6 w-full flex flex-col gap-4'>
			<div class='flex flex-row justify-between'>
				<h2>Ticket Details</h2>
				%s%s
			</div>
			<div class='text-xs'>
				<p><strong>Creator:</strong> %s</p>
				<p><strong>Item:</strong> %s</p>
				<p><strong>Location:</strong> %s</p>
			</div>
			<div class='text-xs'>
				<p><strong>Problem:</strong> %s</p>
			</div>
			<div id='manufacturer-delete-message'>
				%s
			</div>
			<div class='w-[200px]'>
				<img src='data:image/jpeg;base64,%s' class='w-full h-auto' alt='%s'/>
			</div>
			<div id='hidden-settings-section' class='text-xs hidden flex flex-row gap-2'>
				<div id='manufacturer-update-button' class='cursor-pointer rounded py-1 px-2 border border-black bg-white text-black'>Update</div>
				<div id='manufacturer-delete-button' class='cursor-pointer rounded py-1 px-2 border border-black bg-red'>Delete</div>
			</div>
			<div id='update-manufacturer-form' class='hidden mt-2'>
				%s
			</div>
			<div id='delete-manufacturer-form' class='hidden mt-2'>
				%s
			</div>
		</div>
		<script>
			document.getElementById('manufacturer-settings-icon').addEventListener('click', () => {
				document.getElementById('hidden-settings-section').classList.toggle('hidden')
				document.getElementById('manufacturer-close-settings-icon').classList.toggle('hidden')
				document.getElementById('manufacturer-settings-icon').classList.toggle('hidden')
				document.getElementById('manufacturer-delete-message').classList.add('hidden')
			})
			document.getElementById('manufacturer-close-settings-icon').addEventListener('click', () => {
				document.getElementById('hidden-settings-section').classList.toggle('hidden')
				document.getElementById('manufacturer-close-settings-icon').classList.toggle('hidden')
				document.getElementById('manufacturer-settings-icon').classList.toggle('hidden')
				document.getElementById('delete-manufacturer-form').classList.add('hidden')
				document.getElementById('update-manufacturer-form').classList.add('hidden')
			})
			document.getElementById('manufacturer-update-button').addEventListener('click', () => {
				document.getElementById('update-manufacturer-form').classList.toggle('hidden')
				document.getElementById('delete-manufacturer-form').classList.add('hidden')
			})
			document.getElementById('manufacturer-delete-button').addEventListener('click', () => {
				document.getElementById('delete-manufacturer-form').classList.toggle('hidden')
				document.getElementById('update-manufacturer-form').classList.add('hidden')
			})

		</script>
	`,
		SvgIcon("/static/svg/gear-dark.svg", "sm", "id='manufacturer-settings-icon'", ""),
		SvgIcon("/static/svg/x-dark.svg", "sm", "id='manufacturer-close-settings-icon'", "hidden"),
		ticket.Creator,
		ticket.Item,
		ticket.Location,
		ticket.Problem,
		ErrorMessage(err),
		base64.StdEncoding.EncodeToString(ticket.Photo),
		ticket.Item,
		UpdateTicketForm(ticket),
		DeleteTicketForm(ticket),
	)
}

func UpdateTicketForm(ticket _model.Ticket) string {
	return fmt.Sprintf(`
		<form enctype='multipart/form-data' method='POST' action='/app/ticket/%d/update' x-data="{ loading: false }" class='gap-4 grid w-full rounded bg-black'>
			<div class='flex justify-between'>
				%s
			</div>
			<div class='flex flex-col gap-4'>
				%s%s%s%s%s%s
			</div>
		</form>
	`,
		ticket.ID,
		FormTitle("Update Ticket"),
		FormInputLabel("Creator", "creator", "text", ticket.Creator),
		FormInputLabel("Item Description", "item", "text", ticket.Item),
		FormTextAreaLabel("Problem", "problem", 2, ticket.Problem),
		FormPhotoUpload(),
		FormLoader(),
		FormSubmitButton(),
	)
}

func DeleteTicketForm(ticket _model.Ticket) string {
	return fmt.Sprintf(`
		<form x-data="{ loading: false }" action='/app/ticket/%d/delete' method='POST' class='flex flex-col gap-4 w-full'>
			%s%s%s
		</form>
	`, ticket.ID, FormTitle("Delete Ticket"), FormInputLabel("Type 'yes' to delete", "keyword", "", ""), FormDeleteButton())
}

func TicketActivationForm(ticket _model.Ticket, manufacturers []_model.Manufacturer) string {
	manufacturerNames := make([]string, len(manufacturers)+2)
	manufacturerNames[0] = "" // Set the first item as an empty string
	manufacturerNames[1] = "Manufacturer Not Listed"
	for i, manufacturer := range manufacturers {
		manufacturerNames[i+2] = manufacturer.Name
	}
	return fmt.Sprintf(`
		<form id='ticket-activation-form' x-data="{ loading: false }" action='/app/ticket/%d/activation' method='POST' class='flex flex-col gap-4 w-full p-6'>
			%s%s%s%s%s%s
		</form>
		<script>
			let form = document.getElementById('ticket-activation-form')
			console.log(form)
		</script>
	`,
		ticket.ID,
		FormTitle("Ticket Activation"),
		FormInputLabel("Owner", "owner", "text", ticket.Owner),
		FormSelectLabel("Priority", "priority", []string{string(_model.TicketPriorityUnspecified), string(_model.TicketPriorityLow), string(_model.TicketPriorityInconvenient), string(_model.TicketPriorityUrgent)}, string(ticket.Priority)),
		FormSelectLabel("Status", "status", []string{string(_model.TicketStatusNew), string(_model.TicketStatusActive), string(_model.TicketStatusOnHold), string(_model.TicketStatusComplete)}, string(ticket.Status)),
		FormTextAreaLabel("Notes", "notes", 2, ticket.Notes),
		FormSelectLabel("Manufacturer", "manufacturer", manufacturerNames, ""),
	)
}
