import { Theme } from "@radix-ui/themes";
import { createRoot } from "react-dom/client";
import "@radix-ui/themes/styles.css";
import App from "./App.tsx";

createRoot(document.getElementById("root")!).render(
  <Theme>
    <App />
  </Theme>
);
