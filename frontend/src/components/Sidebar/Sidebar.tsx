import { Button, User } from "@gravity-ui/uikit";
import { block } from "../../utils/block";

import "./Sidebar.scss";
import { Link, useLocation } from "react-router";
import axios from "axios";
import { toaster } from "@gravity-ui/uikit/toaster-singleton";
import { useContext, useEffect } from "react";
import { LoginContext } from "../../pages/Login/LoginContext";
import { ButtonWithConfirm } from "../ButtonWithConfirm/ButtonWithConfirm";

const b = block("sidebar");

export const Sidebar = () => {
  const location = useLocation();

  const getLinkClass = (path: string) => {
    return location.pathname === path ? "_active" : "";
  };

  axios.interceptors.response.use(
    (response) => {
      return response;
    },
    (error) => {
      if (error.response.status === 401 || error.response.status === 403) {
        toaster.add({
          title: "Произошла ошибка",
          name: "Отказано в доступе",
          content: "Отказано в доступе",
          theme: "danger",
        });
        return error.response;
      } else if (error.response.status === 400) {
        toaster.add({
          title: "Произошла ошибка",
          name: "Некорректно заполнены поля формы",
          content: "Некорректно заполнены поля формы",
          theme: "danger",
        });
        return error.response;
      } else if (error.response.status !== 404) {
        toaster.add({
          title: "Произошла ошибка",
          name: error.response.data.error,
          content: error.response.data.error,
          theme: "danger",
        });
        return error.response;
      } else {
        return error;
      }
    }
  );

  const loginContext = useContext(LoginContext);

  useEffect(() => {
    const token = localStorage.getItem("token");
    if (token) {
      loginContext.updateUser(token);
    }
  }, []);

  return (
    <div className={b()}>
      {loginContext.userToken ? (
        <>
          <div className={b("user")}>
            <User
              avatar={{ text: loginContext.userName, theme: "brand" }}
              name={loginContext.userName}
              description={loginContext.userRole}
              size="l"
            />
          </div>
          <div className={b("button")}>
            <ButtonWithConfirm
              confirmText="Подтвердите, что хотите выйти из аккаунта"
              handleSubmit={() => {
                loginContext.updateUser("");
              }}
            >
              <Button view="action">Выйти из аккаунта</Button>
            </ButtonWithConfirm>
          </div>

          <div className={b("navigation")}>
            <ul>
              <li>
                <Link to="/" className={getLinkClass("/")}>
                  Главная страница
                </Link>
              </li>
              {loginContext.userRole === "admin" ? (
                <>
                  <li>
                    <Link
                      to="/markup/create"
                      className={getLinkClass("/markup/create")}
                    >
                      Добавить тип разметки
                    </Link>
                  </li>
                  <li>
                    <Link to="/markup" className={getLinkClass("/markup")}>
                      Мои разметки
                    </Link>
                  </li>
                  <li>
                    <Link
                      to="/batch/create"
                      className={getLinkClass("/batch/create")}
                    >
                      Создать проект (batch)
                    </Link>
                  </li>
                  <li>
                    <Link to="/batch" className={getLinkClass("/batch")}>
                      Мои проекты
                    </Link>
                  </li>

                  <li>
                    <Link to="/stat/all" className={getLinkClass("/stat/all")}>
                      Статистика ассессоров
                    </Link>
                  </li>
                </>
              ) : (
                <></>
              )}

              <li>
                <Link to="/stat" className={getLinkClass("/stat")}>
                  Личный кабинет
                </Link>
              </li>
              <li>
                <Link to="/assessment" className={getLinkClass("/assessment")}>
                  Приступить к разметке
                </Link>
              </li>
            </ul>
          </div>
        </>
      ) : (
        <div className={b("buttons")}>
          <Link to="/login" className={b("button")}>
            <Button view="action" width="max">
              Авторизация
            </Button>
          </Link>
          <Link to="/register" className={b("button")}>
            <Button view="action" width="max">
              Регистрация
            </Button>
          </Link>
        </div>
      )}
    </div>
  );
};
