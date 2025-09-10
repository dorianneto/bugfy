import { useState, useEffect } from "react";
import {
  ChevronDownIcon,
  CheckIcon,
  AlertTriangleIcon,
  RefreshCwIcon,
  ExternalLinkIcon,
  CalendarIcon,
  CodeIcon,
  XIcon,
} from "lucide-react";
import { Dialog, Select } from "radix-ui";

// Types
interface Issue {
  id: string;
  project_id: string;
  title: string;
  fingerprint: string;
  count: number;
  first_seen: string;
  last_seen: string;
  status: "unresolved" | "resolved" | "ignored";
}

interface ErrorEvent {
  id: string;
  project_id: string;
  message: string;
  type: string;
  fingerprint: string;
  stack_trace: Array<{
    function: string;
    file: string;
    line: number;
  }>;
  context: Record<string, string>;
  timestamp: string;
  count: number;
}

// API functions
const API_BASE_URL = "http://localhost:8080/api";

const api = {
  getIssues: async (projectId: string): Promise<Issue[]> => {
    const response = await fetch(
      `${API_BASE_URL}/projects/${projectId}/issues`
    );
    if (!response.ok)
      throw new Error(`HTTP ${response.status}: ${response.statusText}`);
    return response.json();
  },

  getErrors: async (fingerprint: string): Promise<ErrorEvent[]> => {
    const response = await fetch(
      `${API_BASE_URL}/issues/${fingerprint}/errors`
    );
    if (!response.ok)
      throw new Error(`HTTP ${response.status}: ${response.statusText}`);
    return response.json();
  },
};

// Utility functions
const formatDistanceToNow = (dateString: string) => {
  try {
    const date = new Date(dateString);
    const now = new Date();
    const diffInSeconds = Math.floor((now.getTime() - date.getTime()) / 1000);

    if (diffInSeconds < 60) return `${diffInSeconds}s ago`;
    if (diffInSeconds < 3600) return `${Math.floor(diffInSeconds / 60)}m ago`;
    if (diffInSeconds < 86400)
      return `${Math.floor(diffInSeconds / 3600)}h ago`;
    return `${Math.floor(diffInSeconds / 86400)}d ago`;
  } catch {
    return "Unknown";
  }
};

const getStatusColor = (status: Issue["status"]) => {
  switch (status) {
    case "resolved":
      return "bg-green-100 text-green-800 border-green-300";
    case "ignored":
      return "bg-gray-100 text-gray-800 border-gray-300";
    default:
      return "bg-red-100 text-red-800 border-red-300";
  }
};

// Main component
export default function Issues() {
  const [issues, setIssues] = useState<Issue[]>([]);
  const [selectedIssue, setSelectedIssue] = useState<Issue | null>(null);
  const [errorDetails, setErrorDetails] = useState<ErrorEvent[]>([]);
  const [loading, setLoading] = useState(true);
  const [loadingDetails, setLoadingDetails] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [projectId, setProjectId] = useState("68b5d94c6be8439becd70248");
  const [statusFilter, setStatusFilter] = useState("all");
  const [searchTerm, setSearchTerm] = useState("");
  const [dialogOpen, setDialogOpen] = useState(false);

  // Load issues
  const loadIssues = async () => {
    if (!projectId.trim()) return;

    try {
      setLoading(true);
      setError(null);
      const data = await api.getIssues(projectId);
      setIssues(data || []);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to load issues");
      setIssues([]);
    } finally {
      setLoading(false);
    }
  };

  // Load error details
  const loadErrorDetails = async (fingerprint: string) => {
    try {
      setLoadingDetails(true);
      const data = await api.getErrors(fingerprint);
      setErrorDetails(data || []);
    } catch (err) {
      setError(
        err instanceof Error ? err.message : "Failed to load error details"
      );
    } finally {
      setLoadingDetails(false);
    }
  };

  // Handle issue click
  const handleIssueClick = async (issue: Issue) => {
    setSelectedIssue(issue);
    setDialogOpen(true);
    await loadErrorDetails(issue.fingerprint);
  };

  useEffect(() => {
    loadIssues();
  }, [projectId]);

  // Filter issues
  const filteredIssues = issues.filter((issue) => {
    const matchesStatus =
      statusFilter === "all" || issue.status === statusFilter;
    const matchesSearch = issue.title
      .toLowerCase()
      .includes(searchTerm.toLowerCase());
    return matchesStatus && matchesSearch;
  });

  return (
    <div className="max-w-6xl mx-auto p-6 space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <h1 className="text-3xl font-bold text-gray-900">Issues</h1>
        <button
          onClick={loadIssues}
          disabled={loading}
          className="inline-flex items-center px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
        >
          <RefreshCwIcon
            className={`w-4 h-4 mr-2 ${loading ? "animate-spin" : ""}`}
          />
          Refresh
        </button>
      </div>

      {/* Filters */}
      <div className="bg-white rounded-lg border border-gray-200 p-6">
        <h3 className="text-lg font-semibold mb-4">Filters</h3>
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          {/* Project ID */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Project ID
            </label>
            <input
              type="text"
              value={projectId}
              onChange={(e) => setProjectId(e.target.value)}
              onKeyDown={(e) => e.key === "Enter" && loadIssues()}
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              placeholder="Enter project ID"
            />
          </div>

          {/* Status Filter */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Status
            </label>
            <Select.Root value={statusFilter} onValueChange={setStatusFilter}>
              <Select.Trigger className="inline-flex items-center justify-between w-full px-3 py-2 text-sm bg-white border border-gray-300 rounded-md hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-blue-500">
                <Select.Value />
                <Select.Icon>
                  <ChevronDownIcon className="w-4 h-4" />
                </Select.Icon>
              </Select.Trigger>
              <Select.Portal>
                <Select.Content className="overflow-hidden bg-white rounded-md shadow-lg border border-gray-200 z-50">
                  <Select.Viewport className="p-1">
                    <Select.Item
                      value="all"
                      className="relative flex items-center px-8 py-2 text-sm rounded hover:bg-gray-100 cursor-default select-none"
                    >
                      <Select.ItemText>All Statuses</Select.ItemText>
                      <Select.ItemIndicator className="absolute left-2">
                        <CheckIcon className="w-4 h-4" />
                      </Select.ItemIndicator>
                    </Select.Item>
                    <Select.Item
                      value="unresolved"
                      className="relative flex items-center px-8 py-2 text-sm rounded hover:bg-gray-100 cursor-default select-none"
                    >
                      <Select.ItemText>Unresolved</Select.ItemText>
                      <Select.ItemIndicator className="absolute left-2">
                        <CheckIcon className="w-4 h-4" />
                      </Select.ItemIndicator>
                    </Select.Item>
                    <Select.Item
                      value="resolved"
                      className="relative flex items-center px-8 py-2 text-sm rounded hover:bg-gray-100 cursor-default select-none"
                    >
                      <Select.ItemText>Resolved</Select.ItemText>
                      <Select.ItemIndicator className="absolute left-2">
                        <CheckIcon className="w-4 h-4" />
                      </Select.ItemIndicator>
                    </Select.Item>
                    <Select.Item
                      value="ignored"
                      className="relative flex items-center px-8 py-2 text-sm rounded hover:bg-gray-100 cursor-default select-none"
                    >
                      <Select.ItemText>Ignored</Select.ItemText>
                      <Select.ItemIndicator className="absolute left-2">
                        <CheckIcon className="w-4 h-4" />
                      </Select.ItemIndicator>
                    </Select.Item>
                  </Select.Viewport>
                </Select.Content>
              </Select.Portal>
            </Select.Root>
          </div>

          {/* Search */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Search
            </label>
            <input
              type="text"
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              placeholder="Search issues..."
            />
          </div>
        </div>
      </div>

      {/* Error State */}
      {error && (
        <div className="bg-red-50 border border-red-200 rounded-lg p-4">
          <div className="flex items-center text-red-600">
            <AlertTriangleIcon className="w-5 h-5 mr-2" />
            <span>{error}</span>
          </div>
        </div>
      )}

      {/* Loading State */}
      {loading && issues.length === 0 ? (
        <div className="flex items-center justify-center py-12">
          <RefreshCwIcon className="w-8 h-8 animate-spin text-blue-600" />
        </div>
      ) : (
        <>
          {/* Issues List */}
          {filteredIssues.length === 0 ? (
            <div className="bg-white rounded-lg border border-gray-200 p-12 text-center">
              <AlertTriangleIcon className="w-12 h-12 mx-auto text-gray-400 mb-4" />
              <h3 className="text-lg font-semibold text-gray-900 mb-2">
                No Issues Found
              </h3>
              <p className="text-gray-500">
                {issues.length === 0
                  ? "No issues have been reported for this project yet."
                  : "No issues match your current filters."}
              </p>
            </div>
          ) : (
            <div className="space-y-4">
              {filteredIssues.map((issue) => (
                <div
                  key={issue.id}
                  onClick={() => handleIssueClick(issue)}
                  className="bg-white rounded-lg border border-gray-200 p-6 hover:shadow-md transition-shadow cursor-pointer"
                >
                  <div className="flex items-start justify-between">
                    <div className="flex-1 min-w-0">
                      <h3 className="text-lg font-semibold text-gray-900 mb-2 break-words">
                        {issue.title}
                      </h3>
                      <div className="flex items-center gap-2 flex-wrap mb-3">
                        <span
                          className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium border ${getStatusColor(
                            issue.status
                          )}`}
                        >
                          {issue.count} events
                        </span>
                        <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-gray-100 text-gray-800 border border-gray-300">
                          {issue.status}
                        </span>
                      </div>
                      <div className="text-sm text-gray-500 space-y-1">
                        <div>
                          First seen: {formatDistanceToNow(issue.first_seen)}
                        </div>
                        <div>
                          Last seen: {formatDistanceToNow(issue.last_seen)}
                        </div>
                        <div className="font-mono text-xs bg-gray-100 px-2 py-1 rounded inline-block">
                          {issue.fingerprint.substring(0, 12)}...
                        </div>
                      </div>
                    </div>
                    <button className="p-2 text-gray-400 hover:text-gray-600 transition-colors">
                      <ExternalLinkIcon className="w-4 h-4" />
                    </button>
                  </div>
                </div>
              ))}

              {/* Summary */}
              <div className="text-center text-sm text-gray-500 py-4">
                Showing {filteredIssues.length} of {issues.length} issues
              </div>
            </div>
          )}
        </>
      )}

      {/* Issue Details Modal */}
      <Dialog.Root open={dialogOpen} onOpenChange={setDialogOpen}>
        <Dialog.Portal>
          <Dialog.Overlay className="fixed inset-0 bg-black bg-opacity-50 z-40" />
          <Dialog.Content className="fixed top-[50%] left-[50%] translate-x-[-50%] translate-y-[-50%] bg-white rounded-lg shadow-xl z-50 w-full max-w-4xl max-h-[90vh] overflow-hidden">
            <div className="flex items-center justify-between p-6 border-b border-gray-200">
              <Dialog.Title className="text-xl font-semibold text-gray-900">
                Issue Details
              </Dialog.Title>
              <Dialog.Close className="p-2 text-gray-400 hover:text-gray-600 transition-colors">
                <XIcon className="w-5 h-5" />
              </Dialog.Close>
            </div>

            <div className="p-6 overflow-y-auto max-h-[calc(90vh-80px)]">
              {selectedIssue && (
                <div className="space-y-6">
                  {/* Issue Summary */}
                  <div>
                    <h2 className="text-lg font-semibold text-red-600 mb-2">
                      {selectedIssue.title}
                    </h2>
                    <div className="flex items-center gap-2 mb-4">
                      <span
                        className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium border ${getStatusColor(
                          selectedIssue.status
                        )}`}
                      >
                        {selectedIssue.count} occurrences
                      </span>
                      <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-gray-100 text-gray-800 border border-gray-300">
                        {selectedIssue.status}
                      </span>
                      <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-blue-100 text-blue-800 border border-blue-300">
                        Project: {selectedIssue.project_id}
                      </span>
                    </div>
                    <div className="text-sm text-gray-600 space-y-1">
                      <div className="font-mono bg-gray-100 px-2 py-1 rounded text-xs">
                        Fingerprint: {selectedIssue.fingerprint}
                      </div>
                    </div>
                  </div>

                  {loadingDetails ? (
                    <div className="flex items-center justify-center py-8">
                      <RefreshCwIcon className="w-6 h-6 animate-spin text-blue-600" />
                    </div>
                  ) : (
                    errorDetails.length > 0 && (
                      <>
                        {/* Stack Trace */}
                        {errorDetails[0].stack_trace &&
                          errorDetails[0].stack_trace.length > 0 && (
                            <div>
                              <h3 className="flex items-center text-lg font-semibold mb-3">
                                <CodeIcon className="w-5 h-5 mr-2" />
                                Stack Trace
                              </h3>
                              <div className="bg-gray-50 rounded-md p-4 space-y-2">
                                {errorDetails[0].stack_trace.map(
                                  (frame, index) => (
                                    <div
                                      key={index}
                                      className="border-l-2 border-gray-300 pl-4"
                                    >
                                      <div className="font-mono text-sm font-semibold text-gray-900">
                                        {frame.function}
                                      </div>
                                      <div className="text-xs text-gray-500">
                                        {frame.file}:{frame.line}
                                      </div>
                                    </div>
                                  )
                                )}
                              </div>
                            </div>
                          )}

                        {/* Context */}
                        {errorDetails[0].context &&
                          Object.keys(errorDetails[0].context).length > 0 && (
                            <div>
                              <h3 className="text-lg font-semibold mb-3">
                                Context
                              </h3>
                              <div className="bg-gray-50 rounded-md p-4 space-y-2">
                                {Object.entries(errorDetails[0].context).map(
                                  ([key, value]) => (
                                    <div key={key} className="flex gap-2">
                                      <span className="font-mono text-sm font-semibold text-gray-700 min-w-0 flex-shrink-0">
                                        {key}:
                                      </span>
                                      <span className="font-mono text-sm text-gray-600 break-all">
                                        {value}
                                      </span>
                                    </div>
                                  )
                                )}
                              </div>
                            </div>
                          )}

                        {/* Recent Events */}
                        <div>
                          <h3 className="flex items-center text-lg font-semibold mb-3">
                            <CalendarIcon className="w-5 h-5 mr-2" />
                            Recent Events ({errorDetails.length})
                          </h3>
                          <div className="space-y-3">
                            {errorDetails.slice(0, 10).map((errorEvent) => (
                              <div
                                key={errorEvent.id}
                                className="border border-gray-200 rounded-md p-4"
                              >
                                <div className="flex items-center justify-between mb-2">
                                  <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-gray-100 text-gray-800 border border-gray-300">
                                    {errorEvent.count} times
                                  </span>
                                  <span className="text-sm text-gray-500">
                                    {new Date(
                                      errorEvent.timestamp
                                    ).toLocaleString()}
                                  </span>
                                </div>
                                {errorEvent.context &&
                                  Object.keys(errorEvent.context).length >
                                    0 && (
                                    <div className="text-xs text-gray-500 space-y-1">
                                      {Object.entries(errorEvent.context).map(
                                        ([key, value]) => (
                                          <div key={key}>
                                            <span className="font-semibold">
                                              {key}:
                                            </span>{" "}
                                            {value}
                                          </div>
                                        )
                                      )}
                                    </div>
                                  )}
                              </div>
                            ))}

                            {errorDetails.length > 10 && (
                              <div className="text-center">
                                <p className="text-sm text-gray-500">
                                  Showing 10 of {errorDetails.length} events
                                </p>
                              </div>
                            )}
                          </div>
                        </div>
                      </>
                    )
                  )}
                </div>
              )}
            </div>
          </Dialog.Content>
        </Dialog.Portal>
      </Dialog.Root>
    </div>
  );
}
