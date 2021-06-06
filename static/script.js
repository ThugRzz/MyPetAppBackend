async function test() {
    let userId = window.location.pathname
    console.log(window.location.pathname);
    let response = await fetch("http://localhost:8000/api/user/qr" + userId);
    if (response.ok) {
        let json = await response.json();
        console.log(json);
        document.getElementById("name").innerHTML = json.data.pet_name;
        document.getElementById("type").innerHTML = json.data.pet;
        document.getElementById("breed").innerHTML = json.data.breed;
        document.getElementById("age").innerHTML = json.data.pet_age;
        document.getElementById("master-name").innerHTML = json.data.owner_name;
        document.getElementById("phone").innerHTML = json.data.phone;
        document.getElementById("email").innerHTML = json.data.email;
        document.getElementById("address").innerHTML = json.data.address;
    } else {
        console.log("Ошибка HTTP: " + response.status);
    }
}

test()