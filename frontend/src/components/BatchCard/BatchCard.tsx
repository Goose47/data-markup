import { Button, Card, Switch, Tooltip } from "@gravity-ui/uikit";
import { block } from "../../utils/block";
import { BatchCardType } from "../../utils/types";
import "./BatchCard.scss";
import { Link } from "react-router";
import { PencilToLine } from "@gravity-ui/icons";
import { useEffect, useRef } from "react";
import { batchUpdate } from "../../utils/requests";

const b = block("batch-card");

const _ = require("lodash");

type BatchCardProps = {
  batch: BatchCardType;
  handleUpdateBatch: (batch: BatchCardType) => void;
  triggerRerender: () => void;
};

export const BatchCard = ({
  batch,
  triggerRerender,
  handleUpdateBatch,
}: BatchCardProps) => {
  const wasUpdated = useRef(false);

  const handleUpdateStatus = (checked: boolean) => {
    const batchCopy = _.cloneDeep(batch);
    batchCopy.is_active = checked;
    handleUpdateBatch(batchCopy);
    wasUpdated.current = true;
  };

  useEffect(() => {
    if (wasUpdated.current) {
      batchUpdate(batch).then(() => {
        "nothing";
      });
    }
    wasUpdated.current = false;
  }, [batch.is_active]);

  return (
    <div className={b()}>
      <Card className={b("card")} onClick={() => triggerRerender()}>
        <div className={b("wrapper")}>
          <div className={b("edit")}>
            <Link to={`/batch/${batch.id}/edit`}>
              <PencilToLine width={20} height={20}></PencilToLine>
            </Link>
          </div>

          <div className={b("toggle")}>
            <div>Проект в работе у ассессоров?</div>
            <Switch
              checked={batch.is_active}
              onUpdate={handleUpdateStatus}
            ></Switch>
          </div>

          <h2>{batch.name}</h2>
          <p>
            Количество пересечений: {batch.overlaps} <br></br>
            Приоритет: {batch.priority} <br></br>
            Проект находится в статусе "
            {batch.is_active ? "активен" : "неактивен"}" <br></br>
            Тип разметки:{" "}
            {batch.type_id === 1
              ? "простой набор полей"
              : "сравнение двух сущностей"}
          </p>
          <Link to={`/batch/${batch.id}`}>
            <Button view="action">Данные разметки</Button>
          </Link>
        </div>
      </Card>
    </div>
  );
};
