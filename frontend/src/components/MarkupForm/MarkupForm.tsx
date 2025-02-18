import { Button, Card, Modal, Select, TextInput } from "@gravity-ui/uikit";
import { block } from "../../utils/block";
import "./MarkupForm.scss";
import { TrashBin, Plus } from "@gravity-ui/icons";
import { ReactNode, useState } from "react";
import { MarkupTypeField } from "../../utils/types";

const _ = require("lodash");

const b = block("markup-form");

export type MarkupTypeForm = {
  type?: string;
  label: string;
  options: string[];
};

export const MarkupForm = ({
  submit,
  markups,
  setMarkups,
  name,
  setName,
  title,
  description,
  buttonText,
}: {
  submit: (result: MarkupTypeField[]) => void;
  markups: MarkupTypeForm[];
  setMarkups: (markups: MarkupTypeForm[]) => void;
  name: string;
  setName: (name: string) => void;
  title: ReactNode;
  description: ReactNode;
  buttonText: ReactNode;
}) => {
  const handleAddNewMarkup = () => {
    const markupsCopy = _.cloneDeep(markups);
    markupsCopy.push({
      type: undefined,
      label: "",
      options: [""],
    });
    setMarkups(markupsCopy);
  };
  const handleDeleteMarkup = (markupIndex: number) => {
    const markupsCopy = _.cloneDeep(markups);
    markupsCopy.splice(markupIndex, 1);
    setMarkups(markupsCopy);
  };

  const handleUpdateMarkup = (
    markupIndex: number,
    newMarkupType: MarkupTypeForm
  ) => {
    const markupsCopy = _.cloneDeep(markups);
    markupsCopy[markupIndex] = newMarkupType;
    setMarkups(markupsCopy);
  };

  const [open, setOpen] = useState(false);

  const handleSubmit = () => {
    const result: MarkupTypeField[] = [];
    markups.forEach((markup, index) => {
      if (markup.type === "5") {
        result.push({
          name: "",
          label: markup.label,
          group_id: index + 1,
          assessment_type_id: parseInt(markup.type ?? "1"),
        });
        return;
      }
      markup.options.forEach((option) => {
        result.push({
          name: option,
          label: markup.label,
          group_id: index + 1,
          assessment_type_id: parseInt(markup.type ?? "1"),
        });
      });
    });
    console.log(result);
    submit(result);
  };

  return (
    <div className={b()}>
      <h1>{title}</h1>
      <p>{description}</p>
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
          <Card key={markupIndex} className={b("input-type")}>
            <div
              className={b("delete")}
              onClick={() => {
                handleDeleteMarkup(markupIndex);
              }}
            >
              <TrashBin width={20} height={20} />
            </div>
            <label htmlFor={`type-${markupIndex}`}>Выберите тип</label>
            <Select
              placeholder="Тип"
              className={b("select")}
              width={"max"}
              size="l"
              id={`type-${markupIndex}`}
              value={[markup.type ?? ""]}
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
            <div className={b("add-option")}>
              <div className={b("question")}>
                <label htmlFor={`markup-question-${markupIndex}`}>Вопрос</label>
                <TextInput
                  placeholder={`Вопрос (можно оставить пустым)`}
                  id={`markup-question-${markupIndex}`}
                  value={markup.label}
                  onUpdate={(value) => {
                    handleUpdateMarkup(markupIndex, {
                      ...markups[markupIndex],
                      label: value,
                    });
                  }}
                />
              </div>
            </div>

            {markup.type && markup.type !== "5" && (
              <div className={b("add-option")}>
                <label>Варианты ответа:</label>
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
                    <TrashBin
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
                  width={15}
                  height={15}
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
          </Card>
        ))}

        <div onClick={handleAddNewMarkup}>
          <Card className={b("input-type")}>
            <Plus width={20} height={20} />
          </Card>
        </div>
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
      <div className={b("input-group")}></div>
    </div>
  );
};
