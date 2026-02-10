// Modified JavaScript for webcam access and sending frames for analysis

let videoStream = null;
let video = null;
let canvas = null;
let canvasContext = null;
let isAnalyzing = false;
let animationFrameId = null;
const analysisApiUrl = '/recognize_face'; // Endpoint in Flask app

// Function to initialize webcam and start analysis
async function startFacialAnalysis() {
    if (isAnalyzing) {
        stopFacialAnalysis();
        return;
    }

    video = document.getElementById('videoFeed');
    canvas = document.getElementById('canvasOverlay');
    const analyzeButton = document.getElementById('analyzeFaceButton');

    if (!video || !canvas) {
        console.error("Video or canvas element not found. Please ensure they exist in your HTML.");
        alert("Error: Missing video or canvas elements in the page. Cannot start analysis.");
        return;
    }

    canvasContext = canvas.getContext('2d');

    try {
        // Access webcam stream
        videoStream = await navigator.mediaDevices.getUserMedia({ video: true });
        video.srcObject = videoStream;
        await video.play();

        // Adjust canvas size to match video feed
        // Use a slight delay to ensure video dimensions are available
        setTimeout(() => {
            canvas.width = video.videoWidth;
            canvas.height = video.videoHeight;
            isAnalyzing = true;
            if (analyzeButton) {
                analyzeButton.textContent = 'Stop Analysis';
            }
            // Start the analysis loop
            requestAnimationFrame(analyzeVideoFrame);
        }, 500); // 500ms delay

    } catch (err) {
        console.error("Error accessing webcam:", err);
        alert("Error accessing webcam. Please grant permission.");
        stopFacialAnalysis(); // Ensure cleanup if error occurs
    }
}

// Function to stop analysis and release webcam
function stopFacialAnalysis() {
    isAnalyzing = false;
    if (videoStream) {
        videoStream.getTracks().forEach(track => track.stop());
        videoStream = null;
    }
    if (video) {
        video.srcObject = null;
    }
    if (animationFrameId) {
        cancelAnimationFrame(animationFrameId);
        animationFrameId = null;
    }
    // Clear canvas
    if (canvasContext && canvas) {
        canvasContext.clearRect(0, 0, canvas.width, canvas.height);
    }
    const analyzeButton = document.getElementById('analyzeFaceButton');
    if (analyzeButton) {
        analyzeButton.textContent = 'Start Analysis';
    }
    console.log("Facial analysis stopped.");
}

// Function to process a video frame
async function analyzeVideoFrame() {
    if (!isAnalyzing || !video || !canvasContext || !canvas) {
        return;
    }

    // Redraw canvas to clear previous frames and draw current frame
    // Ensure canvas dimensions match video dimensions before drawing
    if (video.videoWidth > 0 && video.videoHeight > 0) {
        canvas.width = video.videoWidth;
        canvas.height = video.videoHeight;
        canvasContext.drawImage(video, 0, 0, canvas.width, canvas.height);
    } else {
        // If video dimensions are not yet available, schedule next frame
        if (isAnalyzing) {
            animationFrameId = requestAnimationFrame(analyzeVideoFrame);
        }
        return;
    }

    // Send to backend API if analyzing and a button is available to control it
    if (isAnalyzing) {
        try {
            const base64Image = canvas.toDataURL('image/jpeg', 0.8); // Capture current frame as JPEG

            // Send to backend API
            const response = await fetch(analysisApiUrl, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ image: base64Image }),
            });

            if (!response.ok) {
                const errorData = await response.json().catch(() => ({})); // Catch potential JSON parsing error
                console.error('API Error:', errorData.error || response.statusText);
                // Optionally display error on UI
            } else {
                const data = await response.json();
                if (data.faces && data.faces.length > 0) {
                    drawFaceBoxes(data.faces);
                } else {
                    // No faces detected, clear previous boxes if any
                    clearFaceBoxes();
                }
            }
        } catch (error) {
            console.error('Error sending frame to API:', error);
            // Handle network errors or other fetch issues
        }
    }

    // Schedule the next frame analysis
    if (isAnalyzing) {
        animationFrameId = requestAnimationFrame(analyzeVideoFrame);
    }
}

// Function to draw boxes around detected faces
function drawFaceBoxes(faces) {
    // Clear previous drawings and redraw the current video frame
    canvasContext.clearRect(0, 0, canvas.width, canvas.height);
    canvasContext.drawImage(video, 0, 0, canvas.width, canvas.height);

    canvasContext.strokeStyle = "rgba(0, 255, 0, 0.8)"; // Green bounding box
    canvasContext.lineWidth = 2;

    faces.forEach(face => {
        const { top, right, bottom, left } = face;
        // Draw rectangle for the face
        canvasContext.strokeRect(left, top, right - left, bottom - top);
        
        // Placeholder for drawing names if API returns them
        // canvasContext.fillStyle = "rgba(0, 255, 0, 0.8)";
        // canvasContext.font = "16px Arial";
        // canvasContext.fillText(face.name, left + 6, top - 6);
    });
}

// Function to clear any drawings on the canvas
function clearFaceBoxes() {
    // Clear canvas and redraw the current video frame to remove old boxes
    canvasContext.clearRect(0, 0, canvas.width, canvas.height);
    canvasContext.drawImage(video, 0, 0, canvas.width, canvas.height);
}


// Initialize event listener for the start/stop button
document.addEventListener('DOMContentLoaded', () => {
    const analyzeButton = document.getElementById('analyzeFaceButton');
    if (analyzeButton) {
        analyzeButton.addEventListener('click', startFacialAnalysis);
    } else {
        console.warn("Button with ID 'analyzeFaceButton' not found. Cannot attach start/stop listener.");
    }
});