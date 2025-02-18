import { useContext, useEffect, useState } from "react";
import {
  MarkupForm,
  MarkupTypeForm,
} from "../../components/MarkupForm/MarkupForm";
import { block } from "../../utils/block";
import { MarkupTypeField, MarkupTypeFull } from "../../utils/types";
import "./MarkupEdit.scss";
import { CircleInfoFill } from "@gravity-ui/icons";
import { useNavigate, useParams } from "react-router";
import {
  getDetailedMarkupType,
  handleEditMarkupType,
} from "../../utils/requests";
import { markTypeBackendToFrontend } from "../../utils/adapters";
import { LoginContext } from "../Login/LoginContext";
import { Loader } from "@gravity-ui/uikit";
import { toaster } from "@gravity-ui/uikit/toaster-singleton";

const b = block("markup-edit");

export const MarkupEdit = () => {
  const [markups, setMarkups] = useState<MarkupTypeForm[]>([]);
  const [name, setName] = useState<string>("");

  const params = useParams();

  const navigate = useNavigate();

  const handleEdit = (result: MarkupTypeField[]) => {
    if (params.id) {
      handleEditMarkupType(params.id, {
        name: name,
        fields: result,
      }).then(() => {
        toaster.add({
          title: "Успешно выполнено",
          name: "Успешно отредактирован тип разметки",
          content: "Успешно отредактирован тип разметки",
          theme: "success",
        });
        navigate("/markup");
      });
    }
  };

  useEffect(() => {
    const id = params.id;
    if (!id) return;
    getDetailedMarkupType(parseInt(id)).then((data: MarkupTypeFull) => {
      setMarkups(markTypeBackendToFrontend(data));
      setName(data.name);
    });
  }, [params.id]);

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
        submit={handleEdit}
        title={"Редактирование типа разметки"}
        description={
          <>
            Отредактируйте соответствующую форму и вы сможете добавить этот тип
            разметки к новым проектам.
            <br />
            <br />
            <CircleInfoFill></CircleInfoFill> Тип разметки - это то, что будет
            видеть ассессор в качестве вариантов ответа на определенный запрос в
            проекте (batch-е). Вы в последствие сможете редактировать данный
            шаблон.
            <br />
            <br />
            <CircleInfoFill color="red"></CircleInfoFill> Редактирование данной
            формы не затронет существующие проекты (batch'и). Если вы хотите
            отредактировать тип разметки в существующем проекте - воспользуйтесь
            формой на странице проекта.
          </>
        }
        markups={markups}
        setMarkups={setMarkups}
        name={name}
        setName={setName}
        buttonText={"Редактировать"}
      />
    </div>
  );
};
