import base64
import numpy as np
import cv2
import face_recognition
from flask import Flask, render_template, request, jsonify
import os # Added for potential static file path handling if needed, though Flask handles it.

app = Flask(__name__)

# Route to serve the main dashboard
@app.route('/')
def dashboard():
    return render_template('dashboard.html')

# Route to serve the facial intel panel
@app.route('/facial_intel')
def facial_intel_panel():
    return render_template('facial_intel_panel.html')

# New endpoint for face recognition
@app.route('/recognize_face', methods=['POST'])
def recognize_face():
    try:
        # Get image data from the frontend (expected as base64 string)
        data = request.get_json()
        image_data = data['image'] # Expecting a key named 'image' holding the base64 string

        # Decode the base64 string
        # Remove the header like "data:image/jpeg;base64," if present
        if ',' in image_data:
            image_data = image_data.split(',')[1]
        
        image_bytes = base64.b64decode(image_data)

        # Convert bytes to numpy array for OpenCV
        nparr = np.frombuffer(image_bytes, np.uint8)
        img = cv2.imdecode(nparr, cv2.IMREAD_COLOR)

        if img is None:
            return jsonify({'error': 'Could not decode image'}), 400

        # Convert the image from BGR color (OpenCV default) to RGB color (face_recognition default)
        rgb_img = cv2.cvtColor(img, cv2.COLOR_BGR2RGB)

        # Find all face locations in the image
        face_locations = face_recognition.face_locations(rgb_img)

        # For face recognition (identifying individuals), you would typically:
        # 1. Load known face encodings (e.g., from a directory of images)
        # 2. Encode all faces in the input image
        # 3. Compare encodings to find matches
        # Example:
        # known_face_encodings = [...]
        # known_face_names = [...]
        # face_encodings = face_recognition.face_encodings(rgb_img, face_locations)
        # face_names = []
        # for face_encoding in face_encodings:
        #     matches = face_recognition.compare_faces(known_face_encodings, face_encoding)
        #     name = "Unknown"
        #     if True in matches:
        #         first_match_index = matches.index(True)
        #         name = known_face_names[first_match_index]
        #     face_names.append(name)
        # results = list(zip(face_locations, face_names))

        # For now, we will just return the locations of detected faces.
        # Format locations for JSON: (top, right, bottom, left)
        formatted_locations = []
        for (top, right, bottom, left) in face_locations:
            formatted_locations.append({
                'top': top,
                'right': right,
                'bottom': bottom,
                'left': left
            })

        return jsonify({'faces': formatted_locations})

    except Exception as e:
        print(f"Error processing face recognition: {e}")
        return jsonify({'error': str(e)}), 500

if __name__ == '__main__':
    # Use host='0.0.0.0' to make the server accessible from other devices on the network
    # and port=5000 (default)
    # debug=True is useful during development but should be False in production
    app.run(debug=True, host='0.0.0.0')