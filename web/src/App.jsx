import { useEffect, useState } from "react";
import "./App.css";
import { RoomsPage } from "./pages/RoomsPage.jsx";
import { ChatPage } from "./pages/ChatPage.jsx";

function App() {
  useEffect(() => {
    fetch("http://localhost:8080/ping")
      .then((r) => r.json())
      .then(console.log("Success"));
  }, []);

  return (
    <>
      <ChatPage />
    </>
  );
}

export default App;
