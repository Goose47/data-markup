import { block } from "../../utils/block";
import { BatchCardType } from "../../utils/types";
import "./BatchCard.scss";

const b = block("batch-card");

type BatchCardProps = {
  batch: BatchCardType;
  triggerRerender: () => void;
};

export const BatchCard = ({ batch, triggerRerender }: BatchCardProps) => {
  return (
    <div className={b()} onClick={() => triggerRerender()}>
      {batch.name}
    </div>
  );
};
