import axios from "axios";
import { AssessmentUpdateRq, BatchRq, MarkupTypeRq } from "./types";

const API_PREFIX = "https://api.rwfshr.ru";

export const handleCreateMarkupType = async (request: MarkupTypeRq) => {
  return await axios
    .post(API_PREFIX + "/api/v1/markupTypes", request)
    .then((response) => {
      return response.data;
    });
};

export const handleEditMarkupType = async (
  markupId: string,
  request: MarkupTypeRq
) => {
  return await axios
    .put(API_PREFIX + "/api/v1/markupTypes/" + markupId, request)
    .then((response) => {
      return response.data;
    });
};

export const getAvailableMarkupTypes = async () => {
  return await axios
    .get(API_PREFIX + "/api/v1/markupTypes?batch_id=0")
    .then((response) => response.data.data);
};

export const getDetailedMarkupType = async (markupId: number) => {
  return await axios
    .get(API_PREFIX + "/api/v1/markupTypes/" + markupId)
    .then((response) => response.data);
};

export const deleteMarkupType = async (markupId: number) => {
  return await axios
    .delete(API_PREFIX + "/api/v1/markupTypes/" + markupId)
    .then((response) => response.data);
};

export const createBatch = async (batch: BatchRq) => {
  const form = new FormData();
  form.append("name", batch.name);
  form.append("overlaps", String(batch.overlaps));
  form.append("priority", String(batch.priority));
  form.append("type_id", String(batch.type_id));
  form.append("markups", batch.markups);
  return await axios
    .post(API_PREFIX + "/api/v1/batches/", form)
    .then((response) => response.data);
};

export const linkBatchToMarkupType = async (
  batchId: number,
  markupTypeId: number
) => {
  return await axios
    .post(API_PREFIX + "/api/v1/batches/" + batchId + "/markupTypes", {
      batch_id: batchId,
      markup_type_id: markupTypeId,
      name: "Linked at " + new Date().toString(),
      fields: [],
    })
    .then((response) => response.data);
};

export const getAvailableBatches = async () => {
  return await axios.get(API_PREFIX + "/api/v1/batches").then((response) => {
    return response.data.data;
  });
};

export const assessmentNext = async () => {
  return await axios
    .post(API_PREFIX + "/api/v1/assessments/next")
    .then((response) => response.data);
};

export const assessmentUpdate = async (
  assessmentId: number,
  data: AssessmentUpdateRq
) => {
  return await axios
    .put(API_PREFIX + "/api/v1/assessments/" + assessmentId, data)
    .then((response) => response.data);
};
