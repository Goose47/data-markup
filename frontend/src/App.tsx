import { BrowserRouter, Route, Routes } from "react-router";
import { Home } from "./pages/Home/Home";
import { block } from "./utils/block";

import "./App.scss";
import { Sidebar } from "./components/Sidebar/Sidebar";
import { MyMarkupTypes } from "./pages/MyMarkupTypes/MyMarkupTypes";
import { MarkupCreate } from "./pages/MarkupCreate/MarkupCreate";
import { MarkupEdit } from "./pages/MarkupEdit/MarkupEdit";
import { BatchCreate } from "./pages/BatchCreate/BatchCreate";

const b = block("app");

export const App = () => {
  return (
    <BrowserRouter>
      <div className={b()}>
        <div className={b("wrapper")}>
          <Sidebar />
          <div className={b("content")}>
            <Routes>
              <Route path="/" element={<Home />}></Route>
              <Route path="/markup/create" element={<MarkupCreate />}></Route>
              <Route path="/markup" element={<MyMarkupTypes />}></Route>
              <Route path="/batch/create" element={<BatchCreate />}></Route>
              <Route path="/markup/:id" element={<MarkupEdit />}></Route>
            </Routes>
          </div>
        </div>
      </div>
    </BrowserRouter>
  );
};

export default App;
