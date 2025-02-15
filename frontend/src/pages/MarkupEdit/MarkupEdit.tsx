import { useEffect, useState } from "react";
import {
  MarkupForm,
  MarkupTypeForm,
} from "../../components/MarkupForm/MarkupForm";
import { block } from "../../utils/block";
import { MarkupTypeField, MarkupTypeFull } from "../../utils/types";
import "./MarkupEdit.scss";
import { CircleInfoFill } from "@gravity-ui/icons";
import { useParams } from "react-router";
import {
  getDetailedMarkupType,
  handleEditMarkupType,
} from "../../utils/requests";
import { markTypeBackendToFrontend } from "../../utils/adapters";

const b = block("markup-edit");

export const MarkupEdit = () => {
  const [markups, setMarkups] = useState<MarkupTypeForm[]>([]);
  const [name, setName] = useState<string>("");

  const params = useParams();

  const handleEdit = (result: MarkupTypeField[]) => {
    if (params.id) {
      handleEditMarkupType(params.id, {
        name: name,
        fields: result,
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
