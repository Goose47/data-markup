import { Button, User } from "@gravity-ui/uikit";
import { block } from "../../utils/block";

import "./Sidebar.scss";
import { Link, useLocation, useNavigate } from "react-router";
import axios from "axios";
import { toaster } from "@gravity-ui/uikit/toaster-singleton";

const b = block("sidebar");

export const Sidebar = () => {
  const location = useLocation();

  const getLinkClass = (path: string) => {
    return location.pathname === path ? "_active" : "";
  };

  const navigate = useNavigate();

  axios.interceptors.response.use(
    (response) => {
      return response;
    },
    (error) => {
      if (error.response.status === 401) {
        toaster.add({
          title: "Произошла ошибка",
          name: "Не существует пользователя с такими данными для входа",
          content: "Не существует пользователя с такими данными для входа",
          theme: "danger",
        });
        localStorage.removeItem("token");
        navigate("/login");
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

  axios.interceptors.request.use(function (config) {
    const token =
      "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjozLCJpc3MiOiJtYXJrdXBzIiwiZXhwIjoxNzM5ODkxNzE1LCJpYXQiOjE3Mzk4MDUzMTV9.zXgZgdkas5QTrZcSkfQAnqtAqEyHFGvFBdF6Hv3R4Gg";
    config.headers.Authorization = token;

    return config;
  });

  return (
    <div className={b()}>
      <div className={b("user")}>
        <User
          avatar={{ text: "Charles Darwin", theme: "brand" }}
          name="Charles Darwin"
          description="ассессор"
          size="l"
        />
      </div>
      <div className={b("button")}>
        <Button view="action">Выйти из аккаунта</Button>
      </div>
      <div className={b("navigation")}>
        <ul>
          <li>
            <Link to="/" className={getLinkClass("/")}>
              Главная страница
            </Link>
          </li>
          <li>
            <Link
              to="/markup/create"
              className={getLinkClass("/markup/create")}
            >
              Добавить тип разметки
            </Link>
          </li>
          <li>
            <Link to="/batch/create" className={getLinkClass("/batch/create")}>
              Создать проект (batch)
            </Link>
          </li>
          <li>
            <Link to="/batch" className={getLinkClass("/batch")}>
              Мои проекты
            </Link>
          </li>
          <li>
            <Link to="/markup" className={getLinkClass("/markup")}>
              Мои разметки
            </Link>
          </li>
          <li>
            <Link to="/stat/all" className={getLinkClass("/stat/all")}>
              Статистика ассессоров
            </Link>
          </li>
          <li>
            <Link to="/stat" className={getLinkClass("/stat")}>
              Моя статистика
            </Link>
          </li>
        </ul>
      </div>
    </div>
  );
};
