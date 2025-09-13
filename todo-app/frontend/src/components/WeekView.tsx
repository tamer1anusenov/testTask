import React from 'react';
import { format, startOfWeek, addDays, isSameDay, isToday } from 'date-fns';
import { Calendar, ChevronLeft, ChevronRight } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { TaskItem } from './TaskItem';
import { cn } from '@/lib/utils';
import type { Task } from './TodoApp';

interface WeekViewProps {
  tasks: Task[];
  onToggleTask: (id: string) => void;
  onDeleteTask: (id: string) => void;
  onArchiveTask?: (id: string) => void;
  currentWeek: Date;
  onWeekChange: (date: Date) => void;
}

export const WeekView: React.FC<WeekViewProps> = ({
  tasks,
  onToggleTask,
  onDeleteTask,
  onArchiveTask,
  currentWeek,
  onWeekChange,
}) => {
  const weekStart = startOfWeek(currentWeek, { weekStartsOn: 1 }); // Start on Monday
  const weekDays = Array.from({ length: 7 }, (_, i) => addDays(weekStart, i));

  const getTasksForDay = (day: Date) => {
    return tasks.filter(task => isSameDay(new Date(task.deadline), day));
  };

  const goToPreviousWeek = () => {
    onWeekChange(addDays(currentWeek, -7));
  };

  const goToNextWeek = () => {
    onWeekChange(addDays(currentWeek, 7));
  };

  const goToCurrentWeek = () => {
    onWeekChange(new Date());
  };

  return (
    <div className="space-y-6">
      {/* Week Navigation */}
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-4">
          <div className="flex items-center gap-2">
            <Button
              variant="outline"
              size="sm"
              onClick={goToPreviousWeek}
              className="border-app-border hover:bg-app-surface-hover"
            >
              <ChevronLeft className="h-4 w-4" />
            </Button>
            <Button
              variant="outline"
              size="sm"
              onClick={goToNextWeek}
              className="border-app-border hover:bg-app-surface-hover"
            >
              <ChevronRight className="h-4 w-4" />
            </Button>
          </div>
          <h2 className="text-xl font-semibold">
            {format(weekStart, 'MMM d')} - {format(addDays(weekStart, 6), 'MMM d, yyyy')}
          </h2>
        </div>
        <Button
          variant="outline"
          size="sm"
          onClick={goToCurrentWeek}
          className="border-app-border hover:bg-app-surface-hover"
        >
          <Calendar className="h-4 w-4 mr-2" />
          Today
        </Button>
      </div>

      {/* Week Grid */}
      <div className="grid grid-cols-1 md:grid-cols-7 gap-4">
        {weekDays.map((day, index) => {
          const dayTasks = getTasksForDay(day);
          const isCurrentDay = isToday(day);
          
          return (
            <div
              key={day.toISOString()}
              className={cn(
                "bg-card rounded-lg border transition-all duration-200",
                isCurrentDay 
                  ? "border-primary shadow-md bg-gradient-to-br from-primary/5 to-primary/10" 
                  : "border-app-border hover:border-app-border-hover"
              )}
            >
              {/* Day Header */}
              <div className={cn(
                "p-4 border-b",
                isCurrentDay ? "border-primary/20" : "border-app-border"
              )}>
                <div className="flex items-center justify-between">
                  <div>
                    <div className={cn(
                      "text-sm font-medium",
                      isCurrentDay ? "text-primary" : "text-muted-foreground"
                    )}>
                      {format(day, 'EEEE')}
                    </div>
                    <div className={cn(
                      "text-lg font-semibold",
                      isCurrentDay ? "text-primary" : "text-foreground"
                    )}>
                      {format(day, 'd')}
                    </div>
                  </div>
                  {dayTasks.length > 0 && (
                    <div className={cn(
                      "text-xs px-2 py-1 rounded-full",
                      isCurrentDay 
                        ? "bg-primary/20 text-primary" 
                        : "bg-muted text-muted-foreground"
                    )}>
                      {dayTasks.length} {dayTasks.length === 1 ? 'task' : 'tasks'}
                    </div>
                  )}
                </div>
              </div>

              {/* Tasks for the day */}
              <div className="p-3 space-y-2 min-h-[200px]">
                {dayTasks.length === 0 ? (
                  <div className="flex items-center justify-center h-full text-muted-foreground text-sm">
                    {isCurrentDay ? "No tasks for today" : "No tasks"}
                  </div>
                ) : (
                  dayTasks.map((task, taskIndex) => (
                    <div
                      key={task.id}
                      className="animate-fade-in"
                      style={{ animationDelay: `${taskIndex * 50}ms` }}
                    >
                      <div className="scale-95 transform">
                        <TaskItem
                          task={task}
                          onToggle={onToggleTask}
                          onDelete={onDeleteTask}
                          onArchive={onArchiveTask}
                        />
                      </div>
                    </div>
                  ))
                )}
              </div>
            </div>
          );
        })}
      </div>

      {/* Week Summary */}
      <div className="bg-card rounded-lg p-4 border border-app-border">
        <h3 className="font-medium mb-2">Week Summary</h3>
        <div className="flex flex-wrap gap-4 text-sm text-muted-foreground">
          <div>
            Total tasks: <span className="font-medium text-foreground">{tasks.length}</span>
          </div>
          <div>
            Completed: <span className="font-medium text-task-completed">
              {tasks.filter(t => t.completed).length}
            </span>
          </div>
          <div>
            Active: <span className="font-medium text-task-active">
              {tasks.filter(t => !t.completed).length}
            </span>
          </div>
        </div>
      </div>
    </div>
  );
};