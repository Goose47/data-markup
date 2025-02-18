import { BrowserRouter, Route, Routes } from "react-router";
import { Home } from "./pages/Home/Home";
import { block } from "./utils/block";

import "./App.scss";
import { Sidebar } from "./components/Sidebar/Sidebar";
import { MyMarkupTypes } from "./pages/MyMarkupTypes/MyMarkupTypes";
import { MarkupCreate } from "./pages/MarkupCreate/MarkupCreate";
import { MarkupEdit } from "./pages/MarkupEdit/MarkupEdit";
import { BatchCreate } from "./pages/BatchCreate/BatchCreate";
import { Assessment } from "./pages/Assessment/Assessment";
import { MyBatchesCards } from "./pages/MyBatchesCards/MyBatchesCards";
import { BatchMarkup } from "./pages/BatchMarkup/BatchMarkup";
import { BatchEdit } from "./pages/BatchEdit/BatchEdit";
import { MarkupAssessments } from "./pages/MarkupAssessments/MarkupAssessments";
import { Login } from "./pages/Login/Login";
import { Register } from "./pages/Register/Register";

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
              <Route path="/batch" element={<MyBatchesCards />}></Route>
              <Route path="/batch/:batchId" element={<BatchMarkup />}></Route>
              <Route
                path="/batch/markup/:markupId"
                element={<MarkupAssessments />}
              ></Route>
              <Route
                path="/batch/:batchId/edit"
                element={<BatchEdit />}
              ></Route>
              <Route path="/markup/:id" element={<MarkupEdit />}></Route>
              <Route path="/assessment" element={<Assessment />}></Route>
              <Route path="/login" element={<Login />}></Route>
              <Route path="/register" element={<Register />}></Route>
            </Routes>
          </div>
        </div>
      </div>
    </BrowserRouter>
  );
};

export default App;
