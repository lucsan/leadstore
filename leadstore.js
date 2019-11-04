console.log('leadstore')

let Token = ''

const elementsBuilder = () => {
  let output = document.getElementById('leadstore')

  let feedback = document.createElement('div')
  feedback.id = 'feedback'
  output.append(feedback)
  
  let admin = document.createElement('div')
  admin.id = 'admin'
  output.append(admin)
  
  let name = document.createElement('input')
  name.value = 'admin'
  let pass = document.createElement('input')
  pass.value = 'passy'
  let button = document.createElement('button')
  button.id = 'adminButton'
  button.innerText = 'login'
  admin.append(name)
  admin.append(pass)
  admin.append(button)

  let controls = document.createElement('div')
  controls.id = 'controls'
  output.append(controls)

  let display = document.createElement('div')
  display.id = 'display'
  output.append(display)
  
  button.addEventListener('click', login)
}

elementsBuilder()

const controlsBuilder = () => {
  let eList = document.createElement('button')
  eList.innerText = 'List'
  eList.addEventListener('click', callListAll)
  let eAdd = document.createElement('button')
  eAdd.innerText = 'Add'
  eAdd.addEventListener('click', () => callAddLead(eAddInp.value))
  let eAddInp = document.createElement('input')
  eAddInp.id = 'addLead'
  eAddInp.setAttribute('size', '60%')
  eAddInp.value = `{"first":"Chip", "last":"Chong", "email":"chipo@email.co", "company":"Chongo", "postcode":"CH0", "terms":"false"}`
  let controls = document.getElementById('controls')

  controls.append(eList)
  controls.append(eAdd)
  controls.append(eAddInp)

}

const listAll = (rsp) => {
  let display = document.getElementById('display')
  let allList = document.createElement('div')
  allList.id = 'allList'
  display.append(allList)

  rsp.map(i => {
    const d = JSON.parse(i)
    let e = document.createElement('div')
    let es = document.createElement('input')
    let s = `{"id":"${d.Id}", "first":"${d.FirstName}", "last":"${d.LastName}", "email":"${d.Email}", "company":"${d.Company}", "postcode":"${d.Postcode}", "terms":"${d.AcceptTerms}"}`
    es.value = s
    es.id = `input_${d.Id}`
    es.setAttribute('size', '80%')
    let edate = document.createElement('span')
    edate.innerText = `date ${d.DateCreated}`
    let eupd = document.createElement('button')
    eupd.innerText = 'Update'
    eupd.addEventListener('click', () => callUpdate(d.Id))

    let edel = document.createElement('button')
    edel.innerText = 'Delete' 
    edel.addEventListener('click', () => callDelete(d.Id))
    e.append(es)
    e.append(edate)
    e.append(eupd)
    e.append(edel)

    allList.append(e)
  })

}

const callDelete = (id) => {
  console.log(id)
  let xhr = new XMLHttpRequest()
  xhr.open('DELETE', `http://localhost:3000/api/v1/leads/${id}`, true)

  xhr.setRequestHeader('X-Token', Token)  
  xhr.setRequestHeader('Content-type', 'application/json')
  xhr.onreadystatechange = function() {
    if (this.readyState == 4 && this.status == 200 && this.HEADERS_RECEIVED) {
      console.log(`delete${id}`)
      
      let feedback = document.getElementById('feedback')
      feedback.innerText = 'Confirmed: ' + JSON.parse(this.responseText)
      let display = document.getElementById('display')
      display.innerHTML = ''
      callListAll()
    }
  }
  xhr.send()
  
}

const callUpdate = (id) => {
  let ein = document.getElementById(`input_${id}`)
  callAddLead(ein.value)
}

const callListAll = () => {
  let xhr = new XMLHttpRequest()
  xhr.open('GET', 'http://localhost:3000/api/v1/leads/all', true)

  xhr.setRequestHeader('X-Token', Token)  
  xhr.setRequestHeader('Content-type', 'application/json')
  xhr.onreadystatechange = function() {
    if (this.readyState == 4 && this.status == 200 && this.HEADERS_RECEIVED) {
      console.log('listall')
      
      listAll(JSON.parse(this.responseText))
    }
  }
  xhr.send()
}

const callAddLead = (newLead) => {
  let xhr = new XMLHttpRequest()
  xhr.open('POST', 'http://localhost:3000/api/v1/leads/add', true)
  xhr.setRequestHeader('X-Token', Token)  
  xhr.setRequestHeader('Content-type', 'application/json')
  xhr.onreadystatechange = function() {
    if (this.readyState == 4 && this.status == 200 && this.HEADERS_RECEIVED) {
      let feedback = document.getElementById('feedback')
      feedback.innerText = 'Confirmed: ' + JSON.parse(this.responseText)
      let display = document.getElementById('display')
      display.innerHTML = ''
      callListAll()
    }
  }
console.log(newLead)

  let pnl = JSON.parse(newLead)
  let snl = JSON.stringify(pnl)
  let testLead = JSON.stringify({"id":"7", "first":"Gazo", "last":"Gumby", "email":"gumbo@email.co", "company":"Gunbobo", "postcode":"GU5", "terms":"true"}) 

  xhr.send(snl)
}

function login() {
  console.log('login')
  const admin = document.getElementById('admin')
  const data = admin.getElementsByTagName('input')
  const name = data[0].value
  const pword = data[1].value
 
  let xhr = new XMLHttpRequest()
  xhr.open('POST', 'http://localhost:3000/api/v1/login', true)

  xhr.onreadystatechange = function() {
    if (this.readyState == 4 && this.status == 200 && this.HEADERS_RECEIVED) {
      let headers = this.getAllResponseHeaders()
      loginer(this.responseText, headers)
    }
  }
  xhr.send(`{"name":"${name}", "password": "${pword}"}`)
}

const loginer = (rsp, headers) => {
  const hlist = headersList(headers)
  Token = hlist['x-token']
  let feedback = document.getElementById('feedback')
  feedback.innerText = JSON.parse(rsp)
  controlsBuilder()  
}

const headersList = (headers) => {
  // Convert the header string into an array
  // of individual headers
  var arr = headers.trim().split(/[\r\n]+/)

  // Create a map of header names to values
  var headerMap = {}
  arr.forEach(function (line) {
    var parts = line.split(': ')
    var header = parts.shift()
    var value = parts.join(': ')
    headerMap[header] = value
  })
  return headerMap
}
