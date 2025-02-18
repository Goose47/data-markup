import axios from "axios";
import {
  AssessmentStoreRq,
  AssessmentUpdateRq,
  BatchCardType,
  BatchRq,
  MarkupTypeRq,
  MarkupTypeUpdateRq,
} from "./types";

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

export const handleUpdateMarkupTypeLinked = async (
  request: MarkupTypeUpdateRq
) => {
  return await axios
    .post(
      API_PREFIX + "/api/v1/batches/" + request.batch_id + "/markupTypes",
      request
    )
    .then((response) => {
      return response.data;
    });
};

export const getAvailableMarkupTypes = async (
  batchId: number | undefined = undefined
) => {
  return await axios
    .get(API_PREFIX + "/api/v1/markupTypes?batch_id=" + (batchId ? batchId : 0))
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

export const batchUpdate = async (batch: BatchCardType) => {
  return await axios.put(API_PREFIX + "/api/v1/batches/" + batch.id, {
    ...batch,
    id: undefined,
    created_at: undefined,
  });
};

export const getLinkedMarkupsToBatch = async (
  batchId: number,
  page: number,
  perPage: number
) => {
  return await axios
    .get(
      API_PREFIX +
        "/api/v1/markups?" +
        new URLSearchParams({
          batch_id: String(batchId),
          page: String(page),
          per_page: String(perPage),
        }).toString()
    )
    .then((response) => response.data);
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

export const assessmentStore = async (data: AssessmentStoreRq) => {
  return await axios
    .post(API_PREFIX + "/api/v1/assessments", data)
    .then((response) => response.data);
};

export const batchFind = async (batchId: number) => {
  return await axios
    .get(API_PREFIX + "/api/v1/batches/" + batchId)
    .then((response) => response.data);
};

export const downloadBatchResult = (batchId: number) => {
  axios
    .get(API_PREFIX + "/api/v1/batches/" + batchId + "/export", {
      responseType: "blob",
    })
    .then((blob) => {
      const _url = window.URL.createObjectURL(blob.data);
      window.open(_url, "_blank")?.focus();
    });
};

export const getBatchMarkupData = async (markupId: number) => {
  return await axios
    .get(API_PREFIX + "/api/v1/markups/" + markupId)
    .then((response) => response.data);
};

export const handleLogin = async (email: string, password: string) => {
  return await axios
    .post(API_PREFIX + "/api/v1/auth/login", {
      email: email,
      password: password,
    })
    .then((response) => response.data);
};
