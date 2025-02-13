import { CircleInfoFill } from "@gravity-ui/icons";
import {
  MarkupForm,
  MarkupTypeForm,
} from "../../components/MarkupForm/MarkupForm";
import { block } from "../../utils/block";
import "./MarkupCreate.scss";
import { MarkupTypeField } from "../../utils/types";
import { handleCreateMarkupType } from "../../utils/requests";
import { useState } from "react";

const b = block("markup-create");

export const MarkupCreate = () => {
  const [markups, setMarkups] = useState<MarkupTypeForm[]>([]);
  const [name, setName] = useState<string>("");

  const handleCreate = (result: MarkupTypeField[]) => {
    handleCreateMarkupType({ name: name, fields: result });
  };

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
      />
    </div>
  );
};
