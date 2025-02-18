import { useParams } from "react-router";
import { block } from "../../utils/block";
import "./MarkupAssessments.scss";
import { Table } from "@gravity-ui/uikit";
import { useState } from "react";
import { Assessment } from "../Assessment/Assessment";

const b = block("markup-assessments");

export const MarkupAssessments = () => {
  const params = useParams();

  const markupId = params.markupId ?? "0";

  // const [assessments, setAssessments] = useState<>([])

  return (
    <div className={b()}>
      {/* <Table /> */}

      <Assessment markupId={parseInt(markupId)} isAdmin={true} />
    </div>
  );
};
