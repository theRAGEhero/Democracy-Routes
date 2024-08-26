let ref;

window.onload = () => {
    // const api = new JitsiMeetExternalAPI("8x8.vc", {
    //     roomName: "[sensitive_data]",
    //     parentNode: document.querySelector('#jaas-container'),
    //     // Make sure to include a JWT if you intend to record,
    //     // make outbound calls or use any other premium features!
    //     // jwt: [sensitive_data]
    // });

    // document.querySelector('#info').onclick = () => {
    //     console.log("supported commands:");
    //     console.log(api.getSupportedCommands());

    //     console.log("supported events:");
    //     console.log(api.getSupportedEvents());

    //     console.log("available devices:");
    //     api.getAvailableDevices().then(devices => {
    //         console.log(devices);

    //         ref = devices;
    //     });
    // }

    navigator.mediaDevices.getUserMedia({ audio: true }).then((stream) => {
        const recorder = new MediaRecorder(stream);

        let recorded = [];

        recorder.ondataavailable = (e) => {
            recorded.push(e.data);
        }

        recorder.onstop = () => {
            const blob = new Blob(recorded, { type: "audio/ogg; codecs=opus" });
            recorded = [];

            const audioURL = window.URL.createObjectURL(blob);

            const audio = document.createElement("audio");
            audio.src = audioURL;
            audio.controls = true;

            const container = document.querySelector('#recorded');
            container.appendChild(audio);
        }

        document.querySelector('#start').onclick = () => {
            recorder.start();
        }

        document.querySelector('#stop').onclick = () => {
            recorder.stop();
        }
    });
}
