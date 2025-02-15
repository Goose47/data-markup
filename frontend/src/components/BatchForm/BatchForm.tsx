import { Button, Modal, TextInput } from "@gravity-ui/uikit";
import { block } from "../../utils/block";
import "./BatchForm.scss";
import { ReactNode, useState } from "react";

const b = block("batch-form");

export const BatchForm = ({
  submit,
  name,
  setName,
  title,
  description,
  buttonText,
}: {
  submit: () => void;
  name: string;
  setName: (name: string) => void;
  title: ReactNode;
  description: ReactNode;
  buttonText: ReactNode;
}) => {
  const [open, setOpen] = useState(false);

  const handleSubmit = () => {
    submit();
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
    </div>
  );
};
