import { useState } from "react";
import QRCode from "qrcode";

function App() {
  const [url, setURL] = useState("");
  const [tinyURL, setTinyURL] = useState("");
  const [message, setMessage] = useState("Fill in URL below to get started");
  const [imgQR, setImgQR] = useState<string | null>(null);

  function validURL(str) {
    var pattern = new RegExp(
      "^(https?:\\/\\/)?" + // protocol
        "((([a-z\\d]([a-z\\d-]*[a-z\\d])*)\\.)+[a-z]{2,}|" + // domain name
        "((\\d{1,3}\\.){3}\\d{1,3}))" + // OR ip (v4) address
        "(\\:\\d+)?(\\/[-a-z\\d%_.~+]*)*" + // port and path
        "(\\?[;&a-z\\d%_.~+=-]*)?" + // query string
        "(\\#[-a-z\\d_]*)?$",
      "i"
    ); // fragment locator
    return !!pattern.test(str);
  }

  function generateQR() {
    if (!url) console.log(url);
    QRCode.toDataURL(url).then((dataUrl) => {
      setImgQR(dataUrl);
    });
  }

  function handleSubmit(e) {
    e.preventDefault();
    if (validURL(url)) handleShorten();
    else setMessage("Whoops! URL Not Valid! Enter again 👀");
    async function handleShorten() {
      try {
        const response = await fetch(
          `http://localhost:8080/shorten?url=${encodeURIComponent(url)}`
        );
        const data = await response.text();
        console.log("returned data from API", data);
        setTinyURL(data);
        setMessage(data);
      } catch (e) {
        console.log("Error", e);
      }
    }
  }
  return (
    <>
      <div>
        <h1>Welcome to URL Shortner!</h1>
        <h3 className={`message ${tinyURL ? "success" : "error"}`}>
          {message}
        </h3>
        <form className="app-form">
          <div className="input-wrapper">
            <input
              className="url-input"
              type="url"
              value={url}
              placeholder="Enter URL here"
              onChange={(e) => setURL(e.target.value)}
            />
            <button
              className="submit-btn primary-btn"
              type="submit"
              onClick={handleSubmit}
            >
              Shorten URL
            </button>
          </div>
        </form>
        {tinyURL && (
          <>
            <button
              className="secondary-btn copy-btn"
              onClick={() => navigator.clipboard.writeText(tinyURL)}
            >
              Copy URL
            </button>
            <button className="secondary-btn qr-btn" onClick={generateQR}>
              Generate QR
            </button>
          </>
        )}

        <div>
          {imgQR && (
            <img
              className={`qr-code ${imgQR ? "success" : "error"}`}
              src={imgQR}
            />
          )}
        </div>
      </div>
    </>
  );
}

export default App;
