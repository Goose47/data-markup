import { ReactNode, useState } from "react";
import { block } from "../../utils/block";
import "./ButtonWithConfirm.scss";
import { Button, Modal } from "@gravity-ui/uikit";

const b = block("button-with-confirm");

export const ButtonWithConfirm = ({
  children,
  handleSubmit,
  confirmText,
}: {
  children: ReactNode;
  handleSubmit: () => void;
  confirmText?: string;
}) => {
  const [open, setOpen] = useState(false);

  return (
    <div className={b()}>
      <div onClick={() => setOpen(true)}>{children}</div>
      <Modal
        open={open}
        onClose={() => setOpen(false)}
        contentClassName={b("popup")}
      >
        <h1>Подтвердите действие</h1>
        <p>{confirmText}</p>
        <div className={b("popup-button")}>
          <Button onClick={() => setOpen(false)}>Отмена</Button>
          <Button
            view="action"
            onClick={() => {
              handleSubmit();
              setOpen(false);
            }}
          >
            Подтвердить
          </Button>
        </div>
      </Modal>
    </div>
  );
};
