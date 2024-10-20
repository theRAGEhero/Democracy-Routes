let token = "";

window.onload = async () => {
	let onAuthentication = function(response) {
		token = response.token;
		clearError();
		removeLoginForm();
		loadMeetingsForm();
	};

	loadLoginForm(onAuthentication);

	// const tmp = document.getElementById("login-form-template");
	// const foo = tmp.content.cloneNode(true);
	// const app = document.getElementById("app");
	// app.appendChild(foo);
	//
	// const form = document.querySelector('#login-form');
	//
	// form.addEventListener('submit', async (event) => {
	// 	event.preventDefault();
	// 	const formData = new FormData(form);
	// 	const data = {};
	// 	formData.forEach((value, key) => {
	// 		data[key] = value;
	// 	});
	//
	// 	const msg = document.getElementById('error');
	// 	msg.textContent = "";
	//
	// 	try {
	// 		const response = await login(JSON.stringify(data));
	// 		token = response.Token;
	// 		form.remove();
	// 	} catch (error) {
	// 		msg.textContent = error.message;
	// 	}
	// });

	// const api = new JitsiMeetExternalAPI("8x8.vc", {
	//     roomName: "[sensitive_data]",
	//     parentNode: document.querySelector('#jaas-container'),
	//     // Make sure to include a JWT if you intend to record,
	//     // make outbound calls or use any other premium features!
	//     // jwt: [sensitive_data]
	// });

	// var pc = new RTCPeerConnection();
	//
	// var stream = new MediaStream();
	//
	// document.querySelector('#start').onclick = async () => {
	//     pc = new RTCPeerConnection();
	//
	//     stream = await navigator.mediaDevices.getUserMedia({ audio: true });
	//     stream.getTracks().forEach(track => pc.addTrack(track, stream));
	//
	//     const offer = await pc.createOffer();
	//     await pc.setLocalDescription(offer);
	//
	//     const response = await fetch("http://localhost:8080/offer", {
	//         method: "POST",
	//         mode: "no-cors",
	//         headers: {
	//             "Content-Type": "application/json",
	//         },
	//         body: JSON.stringify(offer),
	//     })
	//
	//     if (!response.ok) {
	//         throw new Error("bad response: " + response.status);
	//     }
	//
	//     const answer = await response.json();
	//     await pc.setRemoteDescription(answer);
	// }
	//
	// document.querySelector('#stop').onclick = () => {
	//     pc.close();
	//
	//     stream.getTracks().forEach(track => track.stop());
	// }
}

function showError(error) {
	const errorMessage = document.getElementById('error');
	errorMessage.textContent = error.message;
}

function clearError() {
	const errorMessage = document.getElementById('error');
	errorMessage.textContent = "";
}

function loadLoginForm(onSuccess) {
	const template = document.getElementById("login-form-template");
	const content = template.content.cloneNode(true);
	const app = document.getElementById("app");
	app.appendChild(content);

	const form = document.getElementById("login-form");
	form.addEventListener('submit', async (event) => {
		event.preventDefault();
		const formData = new FormData(form);
		const data = {};
		formData.forEach((value, key) => {
			data[key] = value;
		});

		try {
			const response = await login(JSON.stringify(data));
			onSuccess(response);
		} catch (error) {
			showError(error);
		}
	});
}

function removeLoginForm() {
	const form = document.getElementById("login-form");
	form.remove();
}

function loadMeetingsForm() {
	const template = document.getElementById("meetings-form-template");
	const content = template.content.cloneNode(true);
	const app = document.getElementById("app");
	app.appendChild(content);

	let onMeetingCreated = function(data) {
		document.getElementById("create-meeting-form")?.remove();
		document.getElementById("no-meetings-placeholder")?.remove();

		const table = document.getElementById("meetings-form-table");
		const row = document.createElement("tr");
		const column = document.createElement("td");

		column.textContent = data.title;
		row.id = data.id;

		row.appendChild(column);
		table.appendChild(row);
	}

	const createButton = document.getElementById("create-meeting-button");
	createButton.addEventListener('click', () => {
		if (!document.getElementById("create-meeting-form")) {
			loadCreateMeetingForm(onMeetingCreated);
		}
    });
}

function loadCreateMeetingForm(onSuccess) {
	const template = document.getElementById("create-meeting-form-template");
	const content = template.content.cloneNode(true);
	const app = document.getElementById("app");
	app.appendChild(content);

	const form = document.getElementById("create-meeting-form");
	form.addEventListener('submit', async (event) => {
		event.preventDefault();
		const formData = new FormData(form);
		const data = {};
		formData.forEach((value, key) => {
			data[key] = value;
		});

		try {
			const response = await createMeeting(JSON.stringify(data));
			onSuccess(response);
		} catch (error) {
			showError(error);
		}
	})

	const cancelButton = document.getElementById("cancel-meeting-creation-button");
	cancelButton.addEventListener('click', () => {
		removeCreateMeetingForm();
	});
}

function removeCreateMeetingForm() {
	const form = document.getElementById("create-meeting-form");
	form.remove();
}

async function login(data) {
	const response = await fetch("http://localhost:8080/login", {
		method: "POST",
		body: data,
	});

	const res = await response.json()

	if (!response.ok) {
		throw new Error(res.error);
	}

	return res;
}

async function createMeeting(data) {
	const response = await fetch("http://localhost:8080/meeting", {
		method: "POST",
		headers: {
			"Authorization": "Bearer " + token
		},
		body: data,
	});

	const res = await response.json()

	if (!response.ok) {
		throw new Error(res.error);
	}

	return res;
}
