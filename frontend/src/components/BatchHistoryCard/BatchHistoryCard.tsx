import { block } from "../../utils/block";
import { MarkupType } from "../../utils/types";
import { MyMarkupType } from "../MyMarkupType/MyMarkupType";
import "./BatchHistoryCard.scss";

const b = block("batch-history-card");

export const BatchHistoryCard = ({ item }: { item: MarkupType }) => {
  return (
    <div className={b()}>
      <div className={b("meta")}>
        {item.markup_count ? (
          <h1>
            {(
              (item.correct_assessment_count / item.assessment_count) *
              100
            ).toFixed(2)}
            %
          </h1>
        ) : null}
        Всего оценок (с учетом пересечений): {item.assessment_count}
        <br />
        Корректных оценок: {item.correct_assessment_count} <br />
        Всего разметок: {item.markup_count} <br />
        Добавлено {new Date(item.created_at).toLocaleString("ru-RU")}
      </div>
      <div className={b("data")}>
        <MyMarkupType
          markupType={item}
          triggerRerender={() => {}}
          isAdmin={false}
        />
      </div>
    </div>
  );
};
