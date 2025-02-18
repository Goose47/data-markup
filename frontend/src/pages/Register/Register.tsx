import { Button, TextInput } from "@gravity-ui/uikit";
import { block } from "../../utils/block";
import "./Register.scss";
import { useContext, useState } from "react";
import { handleLogin, handleRegister } from "../../utils/requests";
import { toaster } from "@gravity-ui/uikit/toaster-singleton";
import { sleep } from "../../utils/utils";
import { useNavigate } from "react-router";
import { LoginContext } from "../Login/LoginContext";

const b = block("login");

export const Register = () => {
  const loginContext = useContext(LoginContext);

  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");

  const navigate = useNavigate();

  const handleSubmit = () => {
    if (confirmPassword !== password) {
      toaster.add({
        title: "Ошибка",
        name: "Пароли не совпадают",
        content: "Пароли не совпадают",
        theme: "danger",
      });
      return;
    }
    handleRegister(email, password).then(async () => {
      handleLogin(email, password).then(async (response) => {
        if (response.token) {
          toaster.add({
            title: "Действие выполнено успешно!",
            name: "Вы успешно зарегистрировались",
            content: "Вы успешно зарегистрировались",
            theme: "success",
          });
          loginContext.updateUser(response.token);
          await sleep(500);
          navigate("/");
        }
      });
    });
  };

  return (
    <div className={b()}>
      <h1>Регистрация</h1>
      <div className={b("input-group")}>
        <label htmlFor="email">E-mail</label>
        <TextInput
          size="l"
          placeholder="E-mail"
          id="email"
          type="email"
          value={email}
          onUpdate={(value) => setEmail(value)}
        />
      </div>

      <div className={b("input-group")}>
        <label htmlFor="password">Пароль</label>
        <TextInput
          size="l"
          placeholder="Пароль"
          id="password"
          type="password"
          value={password}
          onUpdate={(value) => setPassword(value)}
        />
      </div>

      <div className={b("input-group")}>
        <label htmlFor="confirm-password">Повторите пароль</label>
        <TextInput
          size="l"
          placeholder="Пароль"
          id="confirm-password"
          type="password"
          value={confirmPassword}
          onUpdate={(value) => setConfirmPassword(value)}
        />
      </div>

      <div className={b("input-group")}>
        <Button view="action" onClick={handleSubmit}>
          Регистрация
        </Button>
      </div>
      <div className={b("input-group")}></div>
    </div>
  );
};
