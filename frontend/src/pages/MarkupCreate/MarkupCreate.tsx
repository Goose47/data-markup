import { CircleInfoFill } from "@gravity-ui/icons";
import {
  MarkupForm,
  MarkupTypeForm,
} from "../../components/MarkupForm/MarkupForm";
import { block } from "../../utils/block";
import "./MarkupCreate.scss";
import { MarkupTypeField } from "../../utils/types";
import { handleCreateMarkupType } from "../../utils/requests";
import { useContext, useState } from "react";
import { LoginContext } from "../Login/LoginContext";
import { Loader } from "@gravity-ui/uikit";
import { toaster } from "@gravity-ui/uikit/toaster-singleton";
import { useNavigate } from "react-router";

const b = block("markup-create");

export const MarkupCreate = () => {
  const [markups, setMarkups] = useState<MarkupTypeForm[]>([]);
  const [name, setName] = useState<string>("");

  const navigate = useNavigate();

  const handleCreate = (result: MarkupTypeField[]) => {
    handleCreateMarkupType({ name: name, fields: result }).then(() => {
      toaster.add({
        title: "Успешно выполнено",
        name: "Успешно добавлен новый тип разметки",
        content: "Успешно добавлен новый тип разметки",
        theme: "success",
      });
      navigate("/markup");
    });
  };

  const loginContext = useContext(LoginContext);
  if (loginContext.loading) {
    return (
      <div
        style={{
          display: "flex",
          justifyContent: "center",
          alignItems: "center",
          height: 500,
        }}
      >
        <Loader></Loader>
      </div>
    );
  }
  if (loginContext.userRole !== "admin") {
    return (
      <div className={b()}>
        <h1>Отказано в доступе</h1>
      </div>
    );
  }

  return (
    <div className={b()}>
      <MarkupForm
        submit={handleCreate}
        title={"Добавление типа разметки"}
        name={name}
        setName={setName}
        description={
          <>
            Заполните соответствующую форму и вы сможете добавить этот тип
            разметки к новым проектам.
            <br />
            <br />
            <CircleInfoFill></CircleInfoFill> Тип разметки - это то, что будет
            видеть ассессор в качестве вариантов ответа на определенный запрос в
            проекте (batch-е). Вы в последствие сможете редактировать данный
            шаблон.
          </>
        }
        markups={markups}
        setMarkups={setMarkups}
        buttonText={"Добавить"}
      />
    </div>
  );
};
