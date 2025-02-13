import { block } from "../../utils/block";
import "./MarkupCreate.scss";
import { MarkupCreateForm } from "../MarkupCreateForm/MarkupCreateForm";

const b = block("markup-create");

export const MarkupCreate = () => {
  return (
    <div className={b()}>
      <MarkupCreateForm />
    </div>
  );
};
