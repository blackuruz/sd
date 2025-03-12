document.addEventListener("DOMContentLoaded", function () {
    if (!window.wails) {
        alert("Wails tidak terdeteksi! Jalankan aplikasi menggunakan 'wails dev'");
    }
    loadConfigs();
});

// Fungsi Aktivasi Lisensi
async function activateLicense() {
    let key = document.getElementById("licenseKey").value;
    try {
        let response = await window.wails.Invoke("ActivateLicense", key);
        document.getElementById("licenseStatus").innerText = response;
    } catch (err) {
        alert("Error: " + err);
    }
}

// Fungsi Deaktivasi Lisensi
async function deactivateLicense() {
    try {
        let response = await window.wails.Invoke("DeactivateLicense");
        document.getElementById("licenseStatus").innerText = response;
    } catch (err) {
        alert("Error: " + err);
    }
}

// Fungsi Tambah Konfigurasi
async function addConfig() {
    let streamKey = document.getElementById("streamKey").value;
    let channelName = document.getElementById("channelName").value;
    let videoFile = document.getElementById("videoFile").files[0];

    if (!streamKey || !channelName || !videoFile) {
        alert("Harap isi semua field!");
        return;
    }

    let reader = new FileReader();
    reader.readAsDataURL(videoFile);
    reader.onload = async function () {
        let config = {
            streamKey: streamKey,
            channelName: channelName,
            videoFile: reader.result // Menggunakan Base64
        };

        try {
            let response = await window.wails.Invoke("AddConfig", config);
            document.getElementById("configStatus").innerText = response;
            loadConfigs();
        } catch (err) {
            alert("Error: " + err);
        }
    };
}

// Fungsi Mengedit Konfigurasi
async function editConfig() {
    let streamKey = document.getElementById("streamKey").value;
    let channelName = document.getElementById("channelName").value;

    if (!streamKey || !channelName) {
        alert("Harap isi semua field!");
        return;
    }

    let config = {
        streamKey: streamKey,
        channelName: channelName
    };

    try {
        let response = await window.wails.Invoke("EditConfig", config);
        document.getElementById("configStatus").innerText = response;
        loadConfigs();
    } catch (err) {
        alert("Error: " + err);
    }
}

// Fungsi Menghapus Konfigurasi
async function deleteConfig() {
    let streamKey = document.getElementById("streamKey").value;
    if (!streamKey) {
        alert("Masukkan Stream Key yang ingin dihapus!");
        return;
    }

    try {
        let response = await window.wails.Invoke("DeleteConfig", streamKey);
        document.getElementById("configStatus").innerText = response;
        loadConfigs();
    } catch (err) {
        alert("Error: " + err);
    }
}

// Fungsi Memulai Streaming
async function startStream() {
    let streamKey = document.getElementById("streamKey").value;
    if (!streamKey) {
        alert("Masukkan Stream Key yang ingin dijalankan!");
        return;
    }

    try {
        let response = await window.wails.Invoke("StartStream", streamKey);
        document.getElementById("streamStatus").innerText = "Streaming started!";
    } catch (err) {
        alert("Error: " + err);
    }
}

// Fungsi Menghentikan Streaming
async function stopStream() {
    let streamKey = document.getElementById("streamKey").value;
    if (!streamKey) {
        alert("Masukkan Stream Key yang ingin dihentikan!");
        return;
    }

    try {
        let response = await window.wails.Invoke("StopStream", streamKey);
        document.getElementById("streamStatus").innerText = "Streaming stopped!";
    } catch (err) {
        alert("Error: " + err);
    }
}

// Fungsi Memuat Konfigurasi dari Backend
async function loadConfigs() {
    try {
        let configs = await window.wails.Invoke("GetConfigs");
        updateConfigList(configs);
    } catch (err) {
        console.error("Error loading configs:", err);
    }
}

// Fungsi Menampilkan Konfigurasi dalam List
function updateConfigList(configs) {
    let listElement = document.getElementById("configList");
    listElement.innerHTML = "";
    
    configs.forEach(config => {
        let listItem = document.createElement("li");
        listItem.textContent = `${config.channelName} - ${config.streamKey}`;
        listElement.appendChild(listItem);
    });
}
