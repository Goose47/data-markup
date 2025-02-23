import { Button, Loader, Modal, Select, TextInput } from "@gravity-ui/uikit";
import { block } from "../../utils/block";
import "./BatchForm.scss";
import { ReactNode, useContext, useEffect, useState } from "react";
import { CircleInfoFill } from "@gravity-ui/icons";
import { Link } from "react-router";
import { MarkupType } from "../../utils/types";
import { getAvailableMarkupTypes } from "../../utils/requests";
import { LoginContext } from "../../pages/Login/LoginContext";

const b = block("batch-form");

export const BatchForm = ({
  submit,
  title,
  description,
  buttonText,
  states,
}: {
  submit: () => void;
  title: ReactNode;
  description: ReactNode;
  buttonText: ReactNode;
  states: {
    name: string;
    setName: (value: string) => void;
    batchType: string;
    setBatchType: (value: string) => void;
    overlaps: string;
    setOverlaps: (value: string) => void;
    priority: string;
    setPriority: (value: string) => void;
    selectedType: string;
    setSelectedType: (value: string) => void;
    file: File | undefined;
    setFile: (value: File) => void;
  };
}) => {
  const [open, setOpen] = useState(false);

  const {
    name,
    setName,
    batchType,
    setBatchType,
    overlaps,
    setOverlaps,
    priority,
    setPriority,
    selectedType,
    setSelectedType,
    file,
    setFile,
  } = states;

  const handleSubmit = () => {
    submit();
  };

  const [types, setTypes] = useState<MarkupType[]>([]);

  useEffect(() => {
    getAvailableMarkupTypes().then((value: MarkupType[]) => {
      setTypes(value);
      setSelectedType(value?.length ? String(value[0].id) : "");
    });
  }, [setSelectedType, setTypes]);

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
      <h1>{title}</h1>
      <p>{description}</p>
      <div className={b("input-group")}>
        <label htmlFor="title">Название разметки</label>
        <TextInput
          placeholder="Введите название, которое будет видно вам для идентификации разметки"
          id="title"
          value={name}
          onUpdate={(value) => setName(value)}
        />
      </div>

      <div className={b("input-group")}>
        <label htmlFor="name">Тип отображения данных у ассессоров</label>
        <Select
          onUpdate={(value) => setBatchType(value[0])}
          value={[batchType]}
          width={"max"}
          placeholder="Выберите тип"
        >
          <Select.Option value="1">Простой набор полей</Select.Option>
          <Select.Option value="2">Сравнение двух сущностей</Select.Option>
        </Select>
      </div>
      {batchType && (
        <>
          <div>
            <p>
              <CircleInfoFill></CircleInfoFill> Теперь необходимо загрузить .csv
              файл со следующей структурой:<br></br>
              <br></br>
              Файл .csv должен обязательно должен иметь первой строкой заголовок
              полей для их идентификации.
              {batchType === "1" && (
                <>
                  Имена полей должны формироваться по следующему принципу:{" "}
                  <br></br>
                  <b>{"{название_поля}_{тип_данных}"}</b>. Пример корректных
                  названий полей файла:
                  <pre>
                    <code>
                      query_text title_text url_url description_text url_img
                    </code>
                  </pre>
                </>
              )}
              {batchType === "2" && (
                <>
                  Имена полей должны формироваться по следующему принципу:{" "}
                  <br></br>
                  <b>{"{название_поля}{1|2|3|4|5|6|7|8|9}_{тип_данных}"}</b>.
                  Пример корректных названий полей файла:
                  <pre>
                    <code>
                      query1_text title1_text url1_url description1_text
                      url1_img <br></br>
                      query2_text title2_text url2_url description2_text
                      url2_img
                    </code>
                  </pre>
                  В интерфейсе это будет отображено в виде сравнения двух
                  сущностей с одинаковыми параметрами. <br></br>
                </>
              )}
              Поддерживаемые форматы в текущей версии:
              <ul>
                <li>text</li>
                <li>url</li>
                <li>img</li>
              </ul>
            </p>
          </div>

          <div className={b("input-group")}>
            <label htmlFor="file">Загрузите файл (*.csv)</label>
            <input
              type="file"
              name="file"
              id="file"
              accept=".csv"
              lang="ru"
              onChange={(e) => {
                if (!e.target.files || !e.target.files.length) return;
                setFile(e.target.files[0]);
              }}
            />
          </div>

          <div className={b("input-group")}>
            <label htmlFor="overlaps">
              Укажите количество совпадений для выявления правильного ответа
            </label>
            <TextInput
              placeholder="Количество пересечений"
              id="overlaps"
              value={overlaps}
              onUpdate={setOverlaps}
            />
          </div>

          <div className={b("input-group")}>
            <label htmlFor="overlaps">Приоритет</label>
            <Select
              onUpdate={(value) => setPriority(value[0])}
              value={[priority]}
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
            <label htmlFor="overlaps">
              Выберите тип разметки (ответы, предоставляемые для ассессоров)
              <br></br>Добавить свой тип разметки можно на{" "}
              <Link to="/markup/create" target="_blank">
                этой странице
              </Link>
            </label>
            <Select
              onUpdate={(value) => setSelectedType(value[0])}
              value={[selectedType]}
              width={"max"}
              onFocus={() => getAvailableMarkupTypes().then(setTypes)}
            >
              {types.map((type) => {
                return (
                  <Select.Option value={String(type.id)}>
                    {type.name}
                  </Select.Option>
                );
              })}
            </Select>
          </div>

          <div className={b("input-group")}>
            <Button view="action" onClick={() => setOpen(true)}>
              {buttonText}
            </Button>
            <Modal
              open={open}
              onClose={() => setOpen(false)}
              contentClassName={b("popup")}
            >
              <h1>Подтвердите действие</h1>
              <p>
                Проверьте, что вы действительно заполнили все поля корректно
              </p>
              <div className={b("popup-button")}>
                <Button onClick={() => setOpen(false)}>Отмена</Button>
                <Button
                  view="action"
                  onClick={() => {
                    handleSubmit();
                    setOpen(false);
                  }}
                >
                  Сохранить
                </Button>
              </div>
            </Modal>
          </div>
          <div className={b("input-group")}></div>
        </>
      )}
    </div>
  );
};
