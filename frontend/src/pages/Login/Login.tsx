import { Button, TextInput } from "@gravity-ui/uikit";
import { block } from "../../utils/block";
import "./Login.scss";
import { useState } from "react";
import { handleLogin } from "../../utils/requests";
import { toaster } from "@gravity-ui/uikit/toaster-singleton";
import { sleep } from "../../utils/utils";
import { useNavigate } from "react-router";

const b = block("login");

export const Login = () => {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");

  const navigate = useNavigate();

  const handleSubmit = () => {
    handleLogin(email, password).then(async (response) => {
      localStorage.setItem("token", response.token);
      toaster.add({
        title: "Действие выполнено успешно!",
        name: "Вы успешно авторизовались",
        content: "Вы успешно авторизовались",
        theme: "success",
      });
      await sleep(500);
      navigate("/");
    });
  };

  return (
    <div className={b()}>
      <h1>Войти в аккаунт</h1>
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
        <Button view="action" onClick={handleSubmit}>
          Авторизоваться
        </Button>
      </div>
      <div className={b("input-group")}></div>
    </div>
  );
};
