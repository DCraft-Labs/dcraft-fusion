import { createRoot } from "react-dom/client";
import { BrowserRouter } from "react-router-dom";
import { App } from "./ui/App.js";
import { createDemoWorkspace } from "./ui/demo-workspace.js";
import "./styles.css";

const root = document.getElementById("root");

if (root !== null) {
  createRoot(root).render(
    <BrowserRouter>
      <App workspace={createDemoWorkspace()} />
    </BrowserRouter>
  );
}
