import { Button, Loader, Select, TextInput } from "@gravity-ui/uikit";
import { block } from "../../utils/block";
import "./BatchEdit.scss";
import {
  BatchCardType,
  MarkupType,
  MarkupTypeField,
  MarkupTypeFull,
} from "../../utils/types";
import { useEffect, useState } from "react";
import {
  batchFind,
  batchUpdate,
  getAvailableMarkupTypes,
  getDetailedMarkupType,
  handleUpdateMarkupTypeLinked,
} from "../../utils/requests";
import { useParams } from "react-router";
import {
  MarkupForm,
  MarkupTypeForm,
} from "../../components/MarkupForm/MarkupForm";
import { CircleInfoFill } from "@gravity-ui/icons";
import { markTypeBackendToFrontend } from "../../utils/adapters";

const b = block("batch-edit");

export const BatchEdit = () => {
  const params = useParams();
  const [batch, setBatch] = useState<BatchCardType>();

  const [name, setName] = useState("");
  const [markups, setMarkups] = useState<MarkupTypeForm[]>([]);

  useEffect(() => {
    if (params.batchId)
      batchFind(parseInt(params.batchId)).then((data: BatchCardType) =>
        setBatch(data)
      );
    if (params.batchId)
      getAvailableMarkupTypes(parseInt(params.batchId)).then(
        (data: MarkupType[]) => {
          getDetailedMarkupType(data[data.length - 1].id).then(
            (data: MarkupTypeFull) => {
              setMarkups(markTypeBackendToFrontend(data));
            }
          );
        }
      );
  }, [params.batchId]);

  const handleUpdateBatch = (updatedBatch: BatchCardType) => {
    setBatch(updatedBatch);
  };

  const handleSubmit = () => {
    if (batch) batchUpdate(batch);
  };

  const handleUpdateMarkupType = (result: MarkupTypeField[]) => {
    handleUpdateMarkupTypeLinked({
      batch_id: parseInt(params.batchId ?? "-1"),
      fields: result,
      name: name,
    });
  };

  if (!batch) {
    return <Loader></Loader>;
  }

  return (
    <div className={b()}>
      <h1>Редактирование проекта</h1>
      <div className={b("input-group")}>
        <label htmlFor="title">Название проекта</label>
        <TextInput
          placeholder="Введите измененное название"
          id="title"
          value={batch?.name ?? ""}
          onUpdate={(value) =>
            handleUpdateBatch({
              ...batch,
              name: value,
            })
          }
        />
      </div>

      <div className={b("input-group")}>
        <label htmlFor="overlaps">Приоритет</label>
        <Select
          onUpdate={(value) =>
            handleUpdateBatch({
              ...batch,
              priority: parseInt(value[0]),
            })
          }
          value={[String(batch?.priority ?? "")]}
          width={"max"}
        >
          <Select.Option value="2">Незначительный</Select.Option>
          <Select.Option value="4">Ниже среднего</Select.Option>
          <Select.Option value="6">Средний</Select.Option>
          <Select.Option value="8">Выше среднего</Select.Option>
          <Select.Option value="10">Критический</Select.Option>
        </Select>
      </div>

      <div className={b("input-group")}>
        <Button view="action" onClick={handleSubmit}>
          Обновить
        </Button>
      </div>
      <div className={b("input-group")}>
        <br></br>
      </div>

      <MarkupForm
        submit={handleUpdateMarkupType}
        name={name}
        setName={setName}
        markups={markups}
        setMarkups={setMarkups}
        title={"Редактирование типа разметки"}
        description={
          <>
            <CircleInfoFill color="red"></CircleInfoFill> Редактирование типа
            разметки не затронет существующие проекты (batch'и) за исключением
            этого. Данная разметка будет применена к текущему проекту и оценки,
            собираемые ассесорами будут иметь новый вид. <br></br>
            Используйте это редактирование, если вам важно сохранить те ответы,
            которые были по предыдущим типам разметки, используйте эту форму.
            Если вы хотите переоценить все данные, создайте новый проект, а
            текущий поставьте в статус "неактивен".
          </>
        }
        buttonText={"Редактировать"}
      ></MarkupForm>
    </div>
  );
};
