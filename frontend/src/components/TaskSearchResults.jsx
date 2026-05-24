function TaskSearchResults({ results, onClear }) {
  const normalizedResults = Array.isArray(results)
    ? results.map((task) => ({
        id: task.id,
        title: task.title,
        description: task.description,
        priority: task.priority,
        status: task.status,
        createdAt: task.created_at,
        updatedAt: task.updated_at,
      }))
    : [];

  const priorityStyles = {
    high: 'bg-red-100 text-red-800 ring-red-200',
    medium: 'bg-yellow-100 text-yellow-800 ring-yellow-200',
    low: 'bg-green-100 text-green-800 ring-green-200',
  };

  const statusLabelStyles = {
    todo: 'bg-slate-100 text-slate-700 ring-slate-200',
    in_progress: 'bg-blue-100 text-blue-800 ring-blue-200',
    done: 'bg-emerald-100 text-emerald-800 ring-emerald-200',
  };

  return (
    <div className="mt-4 rounded-xl border border-slate-200 bg-slate-50 p-4">
      <div className="mb-3 flex items-center justify-between gap-3">
        <h3 className="font-semibold text-slate-900">Search Results</h3>
        <span className="rounded-full bg-white px-3 py-1 text-xs font-medium text-slate-600 shadow-sm ring-1 ring-slate-200">
          {normalizedResults.length} result{normalizedResults.length === 1 ? '' : 's'}
        </span>
      </div>

      <div className="space-y-3">
        {normalizedResults.map((task) => (
          <article
            key={task.id}
            className="rounded-lg border border-white/70 bg-white p-4 shadow-sm transition hover:shadow-md"
          >
            <div className="flex flex-col gap-3 sm:flex-row sm:items-start sm:justify-between">
              <div className="min-w-0 flex-1">
                <div className="flex flex-wrap items-center gap-2">
                  <h4 className="text-base font-semibold text-slate-900">{task.title}</h4>
                  <span
                    className={`rounded-full px-2.5 py-1 text-xs font-semibold ring-1 ${priorityStyles[task.priority] ?? 'bg-slate-100 text-slate-700 ring-slate-200'}`}
                  >
                    {task.priority}
                  </span>
                  <span
                    className={`rounded-full px-2.5 py-1 text-xs font-semibold ring-1 ${statusLabelStyles[task.status] ?? 'bg-slate-100 text-slate-700 ring-slate-200'}`}
                  >
                    {task.status}
                  </span>
                </div>

                <p className="mt-2 text-sm leading-6 text-slate-600">
                  {task.description?.trim() ? task.description : 'No description provided.'}
                </p>
              </div>
            </div>
          </article>
        ))}
      </div>

      <button
        onClick={onClear}
        className="mt-4 text-sm font-medium text-blue-600 hover:underline"
      >
        Clear Results
      </button>
    </div>
  );
}

export default TaskSearchResults;
