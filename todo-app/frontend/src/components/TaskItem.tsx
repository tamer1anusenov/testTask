import React, { useState } from 'react';
import { Check, Trash2, Clock, Calendar, AlertCircle, Archive } from 'lucide-react';
import { format, isToday, isPast, isTomorrow } from 'date-fns';
import { Button } from '@/components/ui/button';
import { Checkbox } from '@/components/ui/checkbox';
import { cn } from '@/lib/utils';

interface Task {
  id: string;
  title: string;
  completed: boolean;
  createdAt: Date;
  deadline: Date;
  priority: 'low' | 'medium' | 'high';
  archived: boolean;
}

interface TaskItemProps {
  task: Task;
  onToggle: (id: string) => void;
  onDelete: (id: string) => void;
  onArchive?: (id: string) => void;
}

export const TaskItem: React.FC<TaskItemProps> = ({ task, onToggle, onDelete, onArchive }) => {
  const [isDeleting, setIsDeleting] = useState(false);

  const handleDelete = () => {
    setIsDeleting(true);
    setTimeout(() => {
      onDelete(task.id);
    }, 200);
  };

  const handleArchive = () => {
    if (onArchive) {
      onArchive(task.id);
    }
  };

  const getPriorityIcon = (priority: 'low' | 'medium' | 'high') => {
    switch (priority) {
      case 'high': return '🔴';
      case 'medium': return '🟡';
      case 'low': return '🟢';
      default: return '🟡';
    }
  };

  const getPriorityColor = (priority: 'low' | 'medium' | 'high') => {
    switch (priority) {
      case 'high': return 'text-red-500';
      case 'medium': return 'text-yellow-500';
      case 'low': return 'text-green-500';
      default: return 'text-yellow-500';
    }
  };

  const formatTime = (date: Date) => {
    const now = new Date();
    const diff = now.getTime() - date.getTime();
    const minutes = Math.floor(diff / 60000);
    const hours = Math.floor(diff / 3600000);
    const days = Math.floor(diff / 86400000);

    if (days > 0) return `${days}d ago`;
    if (hours > 0) return `${hours}h ago`;
    if (minutes > 0) return `${minutes}m ago`;
    return 'Just now';
  };

  const formatDeadline = (deadline: Date) => {
    if (isToday(deadline)) return 'Today';
    if (isTomorrow(deadline)) return 'Tomorrow';
    if (isPast(deadline) && !isToday(deadline)) return `Overdue`;
    return format(deadline, 'MMM d');
  };

  const getDeadlineColor = (deadline: Date, completed: boolean) => {
    if (completed) return 'text-task-completed';
    if (isPast(deadline) && !isToday(deadline)) return 'text-destructive';
    if (isToday(deadline)) return 'text-task-active';
    return 'text-muted-foreground';
  };

  return (
    <div
      className={cn(
        "group bg-card rounded-lg p-4 border transition-all duration-200 hover:shadow-md",
        task.completed 
          ? "bg-task-completed-bg border-task-completed/20" 
          : "border-app-border hover:border-app-border-hover bg-app-surface hover:bg-app-surface-hover",
        isDeleting && "animate-slide-out"
      )}
    >
      <div className="flex items-center gap-3">
        {/* Checkbox */}
        <Checkbox
          checked={task.completed}
          onCheckedChange={() => onToggle(task.id)}
          className={cn(
            "h-5 w-5 transition-colors",
            task.completed 
              ? "border-task-completed data-[state=checked]:bg-task-completed" 
              : "border-app-border hover:border-primary"
          )}
        />

        {/* Task Content */}
        <div className="flex-1 min-w-0">
          <div
            className={cn(
              "font-medium transition-colors",
              task.completed 
                ? "line-through text-task-completed-text" 
                : "text-card-foreground"
            )}
          >
            {task.title}
          </div>
          <div className="flex items-center gap-4 mt-1">
            {/* Priority */}
            <div className="flex items-center gap-1">
              <span className="text-xs">{getPriorityIcon(task.priority)}</span>
              <span className={cn("text-xs font-medium", getPriorityColor(task.priority))}>
                {task.priority.charAt(0).toUpperCase() + task.priority.slice(1)}
              </span>
            </div>

            {/* Creation Time */}
            <div className="flex items-center gap-1">
              <Clock className="h-3 w-3 text-muted-foreground" />
              <span className="text-xs text-muted-foreground">
                {formatTime(task.createdAt)}
              </span>
            </div>
            
            {/* Deadline */}
            <div className="flex items-center gap-1">
              {isPast(task.deadline) && !isToday(task.deadline) && !task.completed ? (
                <AlertCircle className="h-3 w-3 text-destructive" />
              ) : (
                <Calendar className="h-3 w-3 text-muted-foreground" />
              )}
              <span className={cn("text-xs font-medium", getDeadlineColor(task.deadline, task.completed))}>
                {formatDeadline(task.deadline)}
              </span>
            </div>

            {task.completed && (
              <div className="flex items-center gap-1">
                <Check className="h-3 w-3 text-task-completed" />
                <span className="text-xs text-task-completed">Completed</span>
              </div>
            )}

            {task.archived && (
              <div className="flex items-center gap-1">
                <Archive className="h-3 w-3 text-muted-foreground" />
                <span className="text-xs text-muted-foreground">Archived</span>
              </div>
            )}
          </div>
        </div>

        {/* Actions */}
        <div className="flex items-center gap-2 opacity-0 group-hover:opacity-100 transition-opacity">
          {!task.archived && (
            <Button
              variant="outline"
              size="sm"
              onClick={() => onToggle(task.id)}
              className={cn(
                "h-8 px-3 border-app-border",
                task.completed 
                  ? "hover:bg-app-surface text-muted-foreground" 
                  : "hover:bg-task-completed/10 hover:border-task-completed hover:text-task-completed"
              )}
            >
              {task.completed ? 'Undo' : 'Done'}
            </Button>
          )}
          
          {task.completed && !task.archived && onArchive && (
            <Button
              variant="outline"
              size="sm"
              onClick={handleArchive}
              className="h-8 px-3 border-app-border text-blue-600 hover:bg-blue-600/10 hover:border-blue-600 hover:text-blue-700 transition-colors"
            >
              <Archive className="h-3 w-3" />
            </Button>
          )}
          
          <Button
            variant="outline"
            size="sm"
            onClick={handleDelete}
            className="h-8 px-3 text-destructive hover:bg-destructive/10 hover:border-destructive border-app-border"
          >
            <Trash2 className="h-3 w-3" />
          </Button>
        </div>
      </div>
    </div>
  );
};