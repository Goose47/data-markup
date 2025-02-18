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
import { LoginContext, LoginContextType } from "./pages/Login/LoginContext";
import { useCallback, useMemo, useState } from "react";
import { userMe } from "./utils/requests";
import { UserStatPage } from "./pages/UserStatPage/UserStatPage";
import { UserList } from "./pages/UserList/UserList";

const b = block("app");

export const App = () => {
  const [userName, setUserName] = useState("");
  const [userRole, setUserRole] = useState("");
  const [userToken, setUserToken] = useState("");
  const [loading, setLoading] = useState(true);

  const handleUpdateUser = useCallback((token: string) => {
    if (!token) {
      localStorage.removeItem("token");
      setUserName("");
      setUserRole("");
      setUserToken("");
      setLoading(true);
      return;
    }
    localStorage.setItem("token", token);
    userMe().then((data) => {
      setUserRole(
        data?.user?.roles?.map((el: any) => el.name)?.includes("admin")
          ? "admin"
          : "assessor"
      );
      setUserName(data?.user?.email);
      setLoading(false);
    });
    setUserToken(token);
  }, []);

  const ctxValue = useMemo(
    () => ({
      userName: userName,
      userRole: userRole,
      userToken: userToken,
      loading: loading,
      updateUser: handleUpdateUser,
    }),
    [userName, userRole, userToken, loading, handleUpdateUser]
  );

  return (
    <BrowserRouter>
      <div className={b()}>
        <LoginContext.Provider value={ctxValue}>
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
                <Route path="/stat" element={<UserStatPage />}></Route>
                <Route path="/stat/all" element={<UserList />}></Route>
              </Routes>
            </div>
          </div>
        </LoginContext.Provider>
      </div>
    </BrowserRouter>
  );
};

export default App;
