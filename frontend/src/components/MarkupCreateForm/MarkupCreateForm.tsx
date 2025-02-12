import { Button, Modal, Select, TextInput } from "@gravity-ui/uikit";
import { block } from "../../utils/block";
import "./MarkupCreateForm.scss";
import { Minus, Plus } from "@gravity-ui/icons";
import { useState } from "react";
import { CircleInfoFill } from "@gravity-ui/icons";
import { handleCreateMarkupType } from "../../utils/requests";
import { CreateMarkupTypeFields } from "../../utils/types";

const _ = require("lodash");

const b = block("markup-create-form");

export type MarkupTypeCreate = {
  type?: string;
  options: string[];
};

export const MarkupCreateForm = () => {
  const [name, setName] = useState<string>();
  const [markups, setMarkups] = useState<MarkupTypeCreate[]>([]);

  const handleAddNewMarkup = () => {
    const markupsCopy = markups.slice();
    markupsCopy.push({
      type: undefined,
      options: [""],
    });
    setMarkups(markupsCopy);
  };

  const handleUpdateMarkup = (
    markupIndex: number,
    newMarkupType: MarkupTypeCreate
  ) => {
    const markupsCopy = _.cloneDeep(markups);
    markupsCopy[markupIndex] = newMarkupType;
    setMarkups(markupsCopy);
  };

  const [open, setOpen] = useState(false);

  const handleSubmit = () => {
    if (markups.filter((markup) => markup.type).length !== markups.length) {
      alert("Не указан тип данных в типе разметки");
      return;
    }
    if (!name) {
      alert("Имя не должно быть пустым");
      return;
    }

    const result: CreateMarkupTypeFields[] = [];
    markups.forEach((markup, index) => {
      if (markup.type === "5") {
        result.push({
          name: "",
          group_id: index + 1,
          assessment_type_id: parseInt(markup.type ?? "1"),
        });
        return;
      }
      markup.options.forEach((option) => {
        result.push({
          name: option,
          group_id: index + 1,
          assessment_type_id: parseInt(markup.type ?? "1"),
        });
      });
    });
    handleCreateMarkupType({
      name: name,
      fields: result,
    });
  };

  return (
    <div className={b()}>
      <h1>Добавить тип разметки</h1>
      <p>
        Заполните соответствующую форму и вы сможете добавить этот тип разметки
        к новым проектам. <br />
        <br />
        <CircleInfoFill></CircleInfoFill> Тип разметки - это то, что будет
        видеть ассессор в качестве вариантов ответа на определенный запрос в
        проекте (batch-е). Вы в последствие сможете редактировать данный шаблон.
      </p>
      <div className={b("input-group")}>
        <label htmlFor="name">Название разметки</label>
        <TextInput
          size="l"
          placeholder="Введите название, которое будет видно вам для идентификации разметки"
          id="name"
          value={name}
          onUpdate={(value) => setName(value)}
        />
      </div>

      <div className={b("input-group")}>
        {markups.map((markup, markupIndex) => (
          <div key={markupIndex} className={b("input-type")}>
            <label htmlFor={`type-${markupIndex}`}>Выберите тип</label>
            <Select
              placeholder="Тип"
              className={b("select")}
              width={"max"}
              size="l"
              id={`type-${markupIndex}`}
              onUpdate={(values) => {
                handleUpdateMarkup(markupIndex, {
                  ...markups[markupIndex],
                  type: values[0],
                });
              }}
            >
              <Select.Option value="1">RadioButton</Select.Option>
              <Select.Option value="2">Checkbox</Select.Option>
              <Select.Option value="3">Select</Select.Option>
              <Select.Option value="4">Multiselect</Select.Option>
              <Select.Option value="5">Text</Select.Option>
            </Select>

            {markup.type && markup.type !== "5" && (
              <div className={b("add-option")}>
                <div className={b("add-option-title")}>Варианты ответа:</div>
                {markup.options.map((option, index) => (
                  <div key={index} className={b("add-option-option")}>
                    <TextInput
                      placeholder={`Вариант выбора №${index + 1}`}
                      id="name"
                      value={option}
                      onUpdate={(value) => {
                        const optionsCopy = markups[markupIndex].options;
                        optionsCopy[index] = value;
                        handleUpdateMarkup(markupIndex, {
                          ...markups[markupIndex],
                          options: optionsCopy,
                        });
                      }}
                    />
                    <Minus
                      onClick={() => {
                        const optionsCopy = markups[markupIndex].options;
                        optionsCopy.splice(index, 1);
                        handleUpdateMarkup(markupIndex, {
                          ...markups[markupIndex],
                          options: optionsCopy,
                        });
                      }}
                      width={12}
                      height={12}
                    />
                  </div>
                ))}
                <Plus
                  width={20}
                  height={20}
                  onClick={() => {
                    const optionsCopy = markups[markupIndex].options;
                    optionsCopy.push("");
                    handleUpdateMarkup(markupIndex, {
                      ...markups[markupIndex],
                      options: optionsCopy,
                    });
                  }}
                />
              </div>
            )}
          </div>
        ))}

        <div className={b("input-type")} onClick={handleAddNewMarkup}>
          <Plus width={20} height={20} />
        </div>
      </div>
      <div className={b("input-group")}>
        <Button view="action" onClick={() => setOpen(true)}>
          Добавить
        </Button>
        <Modal
          open={open}
          onClose={() => setOpen(false)}
          contentClassName={b("popup")}
        >
          <h1>Подтвердите действие</h1>
          <p>Проверьте, что вы действительно заполнили все поля корректно</p>
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
    </div>
  );
};
