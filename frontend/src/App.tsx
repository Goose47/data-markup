import { BrowserRouter, Route, Routes } from "react-router";
import { Home } from "./pages/Home/Home";
import { block } from "./utils/block";

import "./App.scss";
import { Sidebar } from "./components/Sidebar/Sidebar";
import { MarkupCreate } from "./components/MarkupCreate/MarkupCreate";

const b = block("app");

export const App = () => {
  return (
    <BrowserRouter>
      <div className={b()}>
        <div className={b("wrapper")}>
          <Sidebar />
          <Routes>
            <Route path="/" element={<Home />}></Route>
            <Route path="/markup/create" element={<MarkupCreate />}></Route>
          </Routes>
        </div>
      </div>
    </BrowserRouter>
  );
};

export default App;
