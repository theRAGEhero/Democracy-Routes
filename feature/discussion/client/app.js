let token = "";

window.onload = async () => {
	const tmp = document.getElementById("login-form-template");
	const foo = tmp.content.cloneNode(true);
	const app = document.getElementById("app");
	app.appendChild(foo);

	const form = document.querySelector('#login-form');

	form.addEventListener('submit', async (event) => {
		event.preventDefault();
		const formData = new FormData(form);
		const data = {};
		formData.forEach((value, key) => {
			data[key] = value;
		});

		const msg = document.getElementById('error');
		msg.textContent = "";

		try {
			const response = await login(JSON.stringify(data));
			token = response.Token;
			form.remove();
		} catch (error) {
			msg.textContent = error.message;
		}
	});

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
