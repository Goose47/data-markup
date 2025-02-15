import { useState } from "react";
import { block } from "../../utils/block";
import "./BatchCreate.scss";
import { BatchForm } from "../../components/BatchForm/BatchForm";
import { CircleInfoFill } from "@gravity-ui/icons";

const b = block("batch-create");

export const BatchCreate = () => {
  const [name, setName] = useState<string>("");

  const handleCreate = () => {};

  return (
    <div className={b()}>
      <BatchForm
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
        buttonText={"Добавить"}
      />
    </div>
  );
};
