# URL Shortener Detection Model
using MLJ

function predict_url_shortening(features)
    # Detect common URL shorteners
    shorteners = ["bit.ly", "tinyurl.com", "goo.gl", "ow.ly", "t.co"]
    domain = features["domain"]
    
    is_shortened = any(shortener in domain for shortener in shorteners)
    return is_shortened ? 0.9 : 0.1
end
